package lauma

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1Key              = "lauma-a1"
	a4Key              = "lauma-a4"
	lunarbloomBonusKey = "lauma-lb-bonus"
)

// A1
func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	c.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase(a1Key, -1),
		Amount: func(ai info.AttackInfo) (float64, bool) {
			// Moonsign < 2 | Bloom, Hyperbloom, Burgeon CD
			if c.getMoonsignLevel() < 2 {
				if ai.AttackTag == attacks.AttackTagBloom ||
					ai.AttackTag == attacks.AttackTagHyperbloom ||
					ai.AttackTag == attacks.AttackTagBurgeon {

					c.AddStatMod(character.StatMod{
						Base:         modifier.NewBase(a1Key, -1),
						AffectedStat: attributes.CD,
						Amount: func() ([]float64, bool) {
							val := make([]float64, attributes.EndStatType)
							val[attributes.CD] = 1 // +100% Crit DMG
							return val, true
						},
					})
					return 0, false
				}
				return 0, false
			}

			// Moonsign >= 2 | Lunar-Bloom CD
			if ai.AttackTag == attacks.AttackTagReactionLunarBloom {

				c.AddStatMod(character.StatMod{
					Base:         modifier.NewBase(a1Key, -1),
					AffectedStat: attributes.CD,
					Amount: func() ([]float64, bool) {
						val := make([]float64, attributes.EndStatType)
						val[attributes.CD] = 0.2 // +20% Crit DMG
						return val, true
					},
				})
				return 0, false
			}
			return 0, false
		},
	})

	c.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase(a1Key, -1),
		Amount: func(ai info.AttackInfo) (float64, bool) {
			// Moonsign < 2 | Bloom, Hyperbloom, Burgeon CR
			if c.getMoonsignLevel() < 2 {
				if ai.AttackTag == attacks.AttackTagBloom ||
					ai.AttackTag == attacks.AttackTagHyperbloom ||
					ai.AttackTag == attacks.AttackTagBurgeon {

					c.AddStatMod(character.StatMod{
						Base:         modifier.NewBase(a1Key, -1),
						AffectedStat: attributes.CD,
						Amount: func() ([]float64, bool) {
							val := make([]float64, attributes.EndStatType)
							val[attributes.CR] = 0.15 // +15% Crit Rate
							return val, true
						},
					})
					return 0, false
				}
				return 0, false
			}

			// Moonsign >= 2 | Lunar-Bloom CR
			if ai.AttackTag == attacks.AttackTagReactionLunarBloom {

				c.AddStatMod(character.StatMod{
					Base:         modifier.NewBase(a1Key, -1),
					AffectedStat: attributes.CD,
					Amount: func() ([]float64, bool) {
						val := make([]float64, attributes.EndStatType)
						val[attributes.CR] = 0.1 // +10% Crit Rate
						return val, true
					},
				})
				return 0, false
			}
			return 0, false
		},
	})
}

// A4
func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	// Skill DMG bonus from EM (max 32%)
	em := c.Stat(attributes.EM)
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag(a4Key, -1),
		Amount: func() ([]float64, bool) {
			val := 0.0004 * em
			if val > 0.32 {
				val = 0.32
			}
			m := make([]float64, attributes.EndStatType)
			m[attributes.DmgP] = val
			return m, true
		},
	})

	// TODO: Charge cooldown reduction from EM (max 20%)
}

// Passive talent
// TODO: Check for any errors, Add LunarBloomEnableKey reactable
func (c *char) lunarbloomInit() {
	c.Core.Flags.Custom[reactable.LunarBloomEnableKey] = 1

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)
		em := c.Stat(attributes.EM)
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectLunarBloom:
		case attacks.AttackTagReactionLunarBloom:
		default:
			return false
		}

		bonus := min(em*0.000175, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("lauma adding lunarbloom base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
		return false
	}, lunarbloomBonusKey)
}
