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

package interactionrequests_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/processing/interactionrequests"
	"github.com/superseriousbusiness/gotosocial/testrig"
)

type AcceptTestSuite struct {
	InteractionRequestsTestSuite
}

func (suite *AcceptTestSuite) TestAccept() {
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	var (
		ctx    = context.Background()
		state  = testStructs.State
		acct   = suite.testAccounts["local_account_2"]
		intReq = suite.testInteractionRequests["admin_account_reply_turtle"]
	)

	// Create int reqs processor.
	p := interactionrequests.New(
		testStructs.Common,
		testStructs.State,
		testStructs.TypeConverter,
	)

	apiApproval, errWithCode := p.Accept(ctx, acct, intReq.ID)
	if errWithCode != nil {
		suite.FailNow(errWithCode.Error())
	}

	// Get db interaction approval.
	dbApproval, err := state.DB.GetInteractionApprovalByID(ctx, apiApproval.ID)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Interacting status
	// should now be approved.
	dbStatus, err := state.DB.GetStatusByURI(ctx, dbApproval.InteractionURI)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.False(*dbStatus.PendingApproval)
	suite.Equal(dbApproval.URI, dbStatus.ApprovedByURI)

	// Wait for a notification
	// for interacting status.
	testrig.WaitFor(func() bool {
		notif, err := state.DB.GetNotification(
			ctx,
			gtsmodel.NotificationMention,
			dbStatus.InReplyToAccountID,
			dbStatus.AccountID,
			dbStatus.ID,
		)
		return notif != nil && err == nil
	})
}

func TestAcceptTestSuite(t *testing.T) {
	suite.Run(t, new(AcceptTestSuite))
}