// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package gtsmodel

import "time"

// DomainPermissionIgnore represents one domain that should be ignored
// when domain permission (ignores) are created from subscriptions.
type DomainPermissionIgnore struct {
	ID                 string               `bun:"type:CHAR(26),pk,nullzero,notnull,unique"`                                       // ID of this item in the database.
	CreatedAt          time.Time            `bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"`                    // Time when this item was created.
	PermissionType     DomainPermissionType `bun:",notnull,unique:domain_permission_ignores_permission_type_domain_uniq"`          // Permission type of the ignore.
	Domain             string               `bun:",nullzero,notnull,unique:domain_permission_ignores_permission_type_domain_uniq"` // Domain to ignore. Eg. 'whatever.com'.
	CreatedByAccountID string               `bun:"type:CHAR(26),nullzero,notnull"`                                                 // Account ID of the creator of this ignore.
	CreatedByAccount   *Account             `bun:"-"`                                                                              // Account corresponding to createdByAccountID.
	PrivateComment     string               `bun:",nullzero"`                                                                      // Private comment on this ignore, viewable to admins.
}
