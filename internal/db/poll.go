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

package db

import (
	"context"

	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
)

type Poll interface {
	// GetPollByID fetches the Poll with given ID from the database.
	GetPollByID(ctx context.Context, id string) (*gtsmodel.Poll, error)

	// GetPollByStatusID fetches the Poll with given status ID column value from the database.
	GetPollByStatusID(ctx context.Context, statusID string) (*gtsmodel.Poll, error)

	// PutPoll puts the given Poll in the database.
	PutPoll(ctx context.Context, poll *gtsmodel.Poll) error

	// DeletePollByID deletes the Poll with given ID from the database.
	DeletePollByID(ctx context.Context, id string) error

	// GetPollVoteByID gets the PollVote with given ID from the database.
	GetPollVoteByID(ctx context.Context, id string) (*gtsmodel.PollVote, error)

	// GetPollVotes fetches all PollVotes in Poll with ID, keyed by account ID, from the database.
	GetPollVotes(ctx context.Context, pollID string) (map[string][]*gtsmodel.PollVote, error)

	// GetPollVotesBy fetches the PollVotes in Poll with ID, by account ID, from the database.
	GetPollVotesBy(ctx context.Context, pollID string, accountID string) ([]*gtsmodel.PollVote, error)

	// CountPollVotes counts all PollVotes in Poll with ID, keyed by option index, from the database.
	CountPollVotes(ctx context.Context, pollID string) ([]int, error)

	// PutPollVotes puts the given PollVotes in the database.
	PutPollVotes(ctx context.Context, vote ...*gtsmodel.PollVote) error

	// DeletePollVotes deletes all PollVotes in Poll with given ID from the database.
	DeletePollVotes(ctx context.Context, pollID string) error

	// DeletePollVotesBy deletes the PollVotes in Poll with ID, by account ID, from the database.
	DeletePollVotesBy(ctx context.Context, pollID string, accountID string) error

	// DeletePollVotesByAccountID deletes all PollVotes in all Polls, by account ID, from the database.
	DeletePollVotesByAccountID(ctx context.Context, accountID string) error
}