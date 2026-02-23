// Description: This file contains the implementation of the attack strategy.
package movementstrategy

// AttackStrategy represents the attack strategy.
// It is used to define the attack strategy.
type AttackStrategy struct {
}

// NewAttackStrategy creates a new AttackStrategy.
// It returns a pointer to the new AttackStrategy.
func NewAttackStrategy() *AttackStrategy {
	return &AttackStrategy{}
}

// Move moves the NPC using the attack strategy.
// It moves the NPC using the attack strategy.
func (a *AttackStrategy) Move() string {
	return "The NPC moves using the attack strategy."
}

