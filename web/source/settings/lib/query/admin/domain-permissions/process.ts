/*
	GoToSocial
	Copyright (C) GoToSocial Authors admin@gotosocial.org
	SPDX-License-Identifier: AGPL-3.0-or-later

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

import {
	ParseConfig as CSVParseConfig,
	parse as csvParse
} from "papaparse";
import { nanoid } from "nanoid";

import { isValidDomainPermission, hasBetterScope } from "../../../domain-permission";
import { gtsApi } from "../../gts-api";

import {
	DomainPermInternalKeys,
	isDomainPerms,
	type DomainPerm,
	type DomainPermsImportForm,
} from "./types";

/**
 * entryProcessor builds up a processing function that can be applied to a
 * DomainPermission entry in order to normalize it before submission to the API.
 * @param formData 
 * @returns 
 */
export function entryProcessor(formData: DomainPermsImportForm): (_entry: DomainPerm) => DomainPerm {
	let processingFuncs: { (_entry: DomainPerm): void; }[] = [];

	// Override each obfuscate entry if necessary.
	if (formData.obfuscate !== undefined) {
		const obfuscateEntry = (entry: DomainPerm) => {
			entry.obfuscate = formData.obfuscate;
		};
		processingFuncs.push(obfuscateEntry);
	}

	// Check whether we need to append or replace
	// private_comment and public_comment.
	["private_comment","public_comment"].forEach((commentType) => {
		let text = formData.commentType?.trim();
		if (!text) {
			return;
		}

		switch(formData[`${commentType}_behavior`]) {
			case "append":
				const appendComment = (entry: DomainPerm) => {
					if (entry.commentType == undefined) {
						entry.commentType = text;
					} else {
						entry.commentType = [entry.commentType, text].join("\n");
					}
				};

				processingFuncs.push(appendComment);
				break;
			case "replace":
				const replaceComment = (entry: DomainPerm) => {
					entry.commentType = text;
				};

				processingFuncs.push(replaceComment);
				break;
		}
	});

	return function process(entry) {
		// Call all the assembled processing functions.
		processingFuncs.forEach((f) => f(entry));

		// Unset all internal processing keys
		// and any undefined keys on this entry.
		Object.entries(entry).forEach(([key, val]) => {
			if (DomainPermInternalKeys.has(key) || val == undefined) {
				delete entry[key];
			}
		});

		return entry;
	};
}

/**
 * Parse the given string of domain permissions and return it as an array.
 * Accepts input as a JSON array string, a CSV, or newline-separated domain names.
 * Will throw an error if input is invalid.
 * @param list 
 * @returns
 * @throws
 */
function parseDomainList(list: string): DomainPerm[] {	
	if (list.startsWith("[")) {
		// Assume JSON array.
		const data = JSON.parse(list);
		if (!isDomainPerms(data)) {
			throw "parsed JSON was not array of DomainPermission";
		}

		return data;
	} else if (list.startsWith("#domain") || list.startsWith("domain,severity")) {
		// Assume Mastodon-style CSV.
		const csvParseCfg: CSVParseConfig = {
			header: true,
			// Remove leading '#' if present.
			transformHeader: (header) => header.startsWith("#") ? header.slice(1) : header,
			skipEmptyLines: true,
			dynamicTyping: true
		};
		
		const { data, errors } = csvParse(list, csvParseCfg);
		if (errors.length > 0) {
			let error = "";
			errors.forEach((err) => {
				error += `${err.message} (line ${err.row})`;
			});
			throw error;
		} 

		if (!isDomainPerms(data)) {
			throw "parsed CSV was not array of DomainPermission";
		}

		return data;
	} else {
		// Fallback: assume newline-separated
		// list of simple domain strings.
		const data: DomainPerm[] = [];
		list.split("\n").forEach((line) => {
			let domain = line.trim();
			let valid = true;

			if (domain.startsWith("http")) {
				try {
					domain = new URL(domain).hostname;
				} catch (e) {
					valid = false;
				}
			}

			if (domain.length > 0) {
				data.push({ domain, valid });
			}
		});

		return data;
	}
}

function deduplicateDomainList(list: DomainPerm[]): DomainPerm[] {
	let domains = new Set();
	return list.filter((entry) => {
		if (domains.has(entry.domain)) {
			return false;
		} else {
			domains.add(entry.domain);
			return true;
		}
	});
}

function validateDomainList(list: DomainPerm[]) {
	list.forEach((entry) => {		
		if (entry.domain.startsWith("*.")) {
			// A domain permission always includes
			// all subdomains, wildcard is meaningless here
			entry.domain = entry.domain.slice(2);
		}

		entry.valid = (entry.valid !== false) && isValidDomainPermission(entry.domain);
		if (entry.valid) {
			entry.suggest = hasBetterScope(entry.domain);
		}
		entry.checked = entry.valid;
	});

	return list;
}

/**
 * useProcessDomainPermissionsMutation uses the RTK Query API without actually
 * hitting the GtS API, it's purely an internal function for our own convenience.
 * 
 * It returns the validated and deduplicated domain permission list.
 */
const useProcessDomainPermissionsMutation = gtsApi.injectEndpoints({
	endpoints: (build) => ({
		processDomainPermissions: build.mutation<DomainPerm[], any>({
			async queryFn(formData, _api, _extraOpts, _fetchWithBQ) {
				if (formData.domains == undefined || formData.domains.length == 0) {
					throw "No domains entered";
				}
				
				// Parse + tidy up the form data.
				const permissions = parseDomainList(formData.domains);
				const deduped = deduplicateDomainList(permissions);
				const validated = validateDomainList(deduped);
				
				validated.forEach((entry) => {
					// Set unique key that stays stable
					// even if domain gets modified by user.
					entry.key = nanoid();
				});
				
				return { data: validated };
			}
		})
	})
}).useProcessDomainPermissionsMutation;

export { useProcessDomainPermissionsMutation };