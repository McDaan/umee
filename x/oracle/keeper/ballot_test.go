package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/umee-network/umee/v2/x/oracle/types"
)

func (s *IntegrationTestSuite) TestBallot_OrganizeBallotByDenom() {
	s.app.OracleKeeper.SetExchangeRate(s.ctx, exchangeRate, sdk.OneDec())
	claimMap := make(map[string]types.Claim)

	// Empty Map
	res := s.app.OracleKeeper.OrganizeBallotByDenom(s.ctx, claimMap)
	s.Require().Equal(make(map[string]types.ExchangeRateBallot), res)

	s.app.OracleKeeper.SetAggregateExchangeRateVote(
		s.ctx, valAddr, types.AggregateExchangeRateVote{
			ExchangeRateTuples: types.ExchangeRateTuples{
				types.ExchangeRateTuple{
					Denom:        "UMEE",
					ExchangeRate: sdk.OneDec(),
				},
			},
			Voter: valAddr.String(),
		},
	)

	claimMap[valAddr.String()] = types.Claim{
		Power:     1,
		Weight:    1,
		WinCount:  1,
		Recipient: valAddr,
	}
	res = s.app.OracleKeeper.OrganizeBallotByDenom(s.ctx, claimMap)
	s.Require().Equal(map[string]types.ExchangeRateBallot{
		"UMEE": {
			types.NewVoteForTally(sdk.OneDec(), "UMEE", valAddr, 1),
		},
	}, res)
}

func (s *IntegrationTestSuite) TestBallot_ClearBallots() {
	prevote := types.AggregateExchangeRatePrevote{
		Hash:        "hash",
		Voter:       addr.String(),
		SubmitBlock: 0,
	}
	s.app.OracleKeeper.SetAggregateExchangeRatePrevote(s.ctx, valAddr, prevote)
	prevoteRes, err := s.app.OracleKeeper.GetAggregateExchangeRatePrevote(s.ctx, valAddr)
	s.Require().NoError(err)
	s.Require().Equal(prevoteRes, prevote)

	var tuples types.ExchangeRateTuples
	tuples = append(tuples, types.ExchangeRateTuple{
		Denom:        "UMEE",
		ExchangeRate: sdk.ZeroDec(),
	})
	vote := types.AggregateExchangeRateVote{
		ExchangeRateTuples: tuples,
		Voter:              addr.String(),
	}
	s.app.OracleKeeper.SetAggregateExchangeRateVote(s.ctx, valAddr, vote)
	voteRes, err := s.app.OracleKeeper.GetAggregateExchangeRateVote(s.ctx, valAddr)
	s.Require().NoError(err)
	s.Require().Equal(voteRes, vote)

	s.app.OracleKeeper.ClearBallots(s.ctx, 0)
	_, err = s.app.OracleKeeper.GetAggregateExchangeRatePrevote(s.ctx, valAddr)
	s.Require().Error(err)
	_, err = s.app.OracleKeeper.GetAggregateExchangeRateVote(s.ctx, valAddr)
	s.Require().Error(err)
}
