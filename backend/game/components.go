package game

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type MovementComponent struct {
	MovementSpeed         int
	TimeSinceLastMovement int
	Path                  []engine.Vector
	Target                engine.Vector
}

func (m MovementComponent) ComponentName() string {
	return "MovementComponent"
}

func MovementComponentMaker(data interface{}) engine.Component {
	compData, ok := data.(map[string]interface{})
	if !ok {
		panic(fmt.Sprint("Data given to component isn't correct type, expected map[string]interface{}"))
	}

	component := MovementComponent{}
	component.Target = engine.Vector{X: -1, Y: -1}
	speed, ok := compData[constants.MovementSpeedJson].(string)
	if !ok {
		speed = "slow"
	}

	switch speed {
	case "slow":
		component.MovementSpeed = constants.MovementSpeedSlow
	case "fast":
		component.MovementSpeed = constants.MovementSpeedFast
	}

	return component
}

type PositionComponent struct {
	Position engine.Vector
}

func (m PositionComponent) ComponentName() string {
	return "PositionComponent"
}

func PositionComponentMaker(data interface{}) engine.Component {
	return PositionComponent{}
}

type AttackComponent struct {
	Target              int
	Damage              int
	AttackSpeed         int
	TimeSinceLastAttack int
	Range               int
}

func (a AttackComponent) ComponentName() string {
	return "AttackComponent"
}

func AttackComponentMaker(data interface{}) engine.Component {

	compData, ok := data.(map[string]interface{})

	if !ok {
		panic(fmt.Sprint("Data given to attack component isn't correct type, expected map[string]interface{}"))
	}

	component := AttackComponent{Target: -1}

	attackSpeed := int(compData[constants.AttackSpeedJson].(float64))
	damage := int(compData[constants.DamageJson].(float64))

	component.Damage = damage
	component.AttackSpeed = attackSpeed

	return component
}

type HealthComponent struct {
	MaxHealth int
	Health    int
}

func (h HealthComponent) ComponentName() string {
	return "HealthComponent"
}

func HealthComponentMaker(data interface{}) engine.Component {

	compData, ok := data.(map[string]interface{})

	if !ok {
		panic(fmt.Sprint("Data given to health component isn't correct type, expected map[string]interface{}"))
	}

	component := HealthComponent{}

	health := int(compData[constants.MaxHealthJson].(float64))

	component.MaxHealth = health
	component.Health = health

	return component
}
