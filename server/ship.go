package main

import (
	"fmt"
	"log"
	"math"
)

type Ship struct {
	Id           string  `json:"i"`
	Alive        bool    `json:"z"`
	Position     *Point  `json:"p"`
	Angle        float64 `json:"a"`
	Velocity     *Vector `json:"v"`
	Acceleration int8    `json:"l"`
	Rotation     int8    `json:"r"`
}

const (
	ShipRadius = 10
)

func CreateShip(id string, pos *Point) *Ship {
	return &Ship{
		Id:           id,
		Alive:        true,
		Position:     pos,
		Angle:        0,
		Velocity:     &Vector{0, 0},
		Acceleration: 0,
		Rotation:     0,
	}
}

func (s *Ship) Tick(t uint64, state *GameState) *Ship {
	// calculate time since last update (in milliseconds)
	elapsedMillis := t - state.Time

	// elapsed time in percentage of a second
	elapsed := float64(elapsedMillis) / 1000

	// calculate new angle
	angle := s.Angle
	if s.Rotation != 0 {
		angle = s.Angle + (settings.constants.ShipRotation * elapsed * float64(s.Rotation))
		for angle < 0 {
			angle += 360
		}
		for angle >= 360 {
			angle -= 360
		}
		angle = RoundToPlaces(angle, 1)
	}

	// calculate new velocity
	vel := s.Velocity
	if s.Acceleration == 1 {
		accel := AngleAndSpeedToVector(angle, settings.constants.ShipAcceleration)
		vel = AddVectors(s.Velocity, MultiplyVector(accel, elapsed))
	}

	// apply drag
	vel = AddVectors(vel, MultiplyVector(s.Velocity, settings.constants.ShipDrag*elapsed))

	// calculate new position
	pos := MakePoint(s.Position.X+vel.X*elapsed, s.Position.Y+vel.Y*elapsed)

	// TODO: Come up with a better way to look up collisions.
	// From https://developer.mozilla.org/en-US/docs/Games/Techniques/2D_collision_detection
	for _, os := range state.Ships {
		if os.Id != s.Id {
			dx := s.Position.X - os.Position.X
			dy := s.Position.Y - os.Position.Y
			distance := math.Sqrt(float64(dx*dx + dy*dy))

			if distance < ShipRadius*2 {
				if settings.debug {
					log.Println(fmt.Sprintf("Ship %v colliding with Ship %v", s.Id, os.Id))
				}
			}
		}
	}

	for _, p := range state.Projectiles {
		dx := s.Position.X - p.Position.X
		dy := s.Position.Y - p.Position.Y
		distance := math.Sqrt(float64(dx*dx + dy*dy))

		if distance < ShipRadius+ProjectileRadius {
			if settings.debug {
				log.Println(fmt.Sprintf("Ship %v colliding with Projectile %v", s.Id, p.Id))
			}
		}
	}

	return &Ship{
		Id:           s.Id,
		Alive:        s.Alive,
		Position:     pos,
		Angle:        angle,
		Velocity:     vel,
		Acceleration: s.Acceleration,
		Rotation:     s.Rotation,
	}
}
