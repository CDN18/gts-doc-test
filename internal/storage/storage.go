/*
   GoToSocial
   Copyright (C) 2021 GoToSocial Authors admin@gotosocial.org

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

// Package storage contains an interface and implementations for storing and retrieving files and attachments.
package storage

// Storage is an interface for storing and retrieving blobs
// such as images, videos, and any other attachments/documents
// that shouldn't be stored in a database.
type Storage interface {
	StoreFileAt(path string, data []byte) error
	RetrieveFileFrom(path string) ([]byte, error)
}
