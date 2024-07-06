package nysenateapi

type VoteEntry struct {
	MemberID  int
	Chamber   string
	VoteType  string // COMMITTEE, FLOOR
	Vote      string // Aye, Nay, Excused
	ShortName string
}
type VoteEntries []VoteEntry

func (v VoteEntries) Filter(chamber string) VoteEntries {
	var o VoteEntries
	for _, vv := range v {
		if vv.Chamber == chamber {
			o = append(o, vv)
		}
	}
	return o
}

func (b Bill) GetVotes() VoteEntries {
	var o VoteEntries
	// TODO: dedupe
	for _, v := range b.Votes.Items {
		if v.Version != b.ActiveVersion && v.Version != "" {
			// Assembly workaround Votes don't have version
			// TODO: all versions?
			continue
		}
		for _, m := range v.MemberVotes.Items.Excused.Items {
			o = append(o, VoteEntry{
				ShortName: m.ShortName,
				MemberID:  m.MemberID,
				Chamber:   m.Chamber,
				VoteType:  v.VoteType,
				Vote:      "Excused",
			})
		}
		for _, m := range v.MemberVotes.Items.Aye.Items {
			o = append(o, VoteEntry{
				ShortName: m.ShortName,
				MemberID:  m.MemberID,
				Chamber:   m.Chamber,
				VoteType:  v.VoteType,
				Vote:      "Aye",
			})
		}
		for _, m := range v.MemberVotes.Items.Nay.Items {
			o = append(o, VoteEntry{
				ShortName: m.ShortName,
				MemberID:  m.MemberID,
				Chamber:   m.Chamber,
				VoteType:  v.VoteType,
				Vote:      "Nay",
			})
		}
		for _, m := range v.MemberVotes.Items.AyeWithReservations.Items {
			o = append(o, VoteEntry{
				ShortName: m.ShortName,
				MemberID:  m.MemberID,
				Chamber:   m.Chamber,
				VoteType:  v.VoteType,
				Vote:      "Aye",
				// TODO: add note "with reservations"
			})
		}
		for _, m := range v.MemberVotes.Items.Absent.Items {
			o = append(o, VoteEntry{
				ShortName: m.ShortName,
				MemberID:  m.MemberID,
				Chamber:   m.Chamber,
				VoteType:  v.VoteType,
				Vote:      "Absent",
			})
		}
		// TODO: Abstained ?
	}
	return o
}
