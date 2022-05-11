package deadball

type Trait struct {
	Label string
	Name  string
}

var (
	TraitPowerPlus Trait = Trait{
		Label: "P++",
		Name:  "Power Hitter +",
	}
	TraitPower Trait = Trait{
		Label: "P+",
		Name:  "Power Hitter",
	}
	TraitContact Trait = Trait{
		Label: "C+",
		Name:  "Contact Hitter",
	}
	TraitSpeedy Trait = Trait{
		Label: "S+",
		Name:  "Speedy Runner",
	}
	TraitDefender Trait = Trait{
		Label: "D+",
		Name:  "Great Defender",
	}
	TraitWeakMinus Trait = Trait{
		Label: "P--",
		Name:  "Weak Hitter -",
	}
	TraitWeak Trait = Trait{
		Label: "P-",
		Name:  "Weak Hitter",
	}
	TraitFreeSwinger Trait = Trait{
		Label: "C-",
		Name:  "Free Swinger",
	}
	TraitSlow Trait = Trait{
		Label: "S-",
		Name:  "Slow Runner",
	}
	TraitPoorDefender Trait = Trait{
		Label: "D-",
		Name:  "Poor Defender",
	}
)
