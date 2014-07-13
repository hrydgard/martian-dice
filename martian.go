// Copyright (c) 2014 Henrik Rydgård

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2.0 or later versions.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License 2.0 for more details.

// See http://www.gnu.org/licenses/ for a full copy of the license.

package main

import (
  "fmt"
  "bufio"
  "os"
  "strings"
  "time"
  "math/rand"
  "errors"
)

type Die struct {
  value int
  locked bool
}

const (
  tank = iota
  deathray
  human
  cow
  chicken
)

type PlayerState struct {
  Score int
}

type MartianState struct {
  // Active Dice
  dice [13]Die

  keptTanks int
  keptDeathrays int
  keptHumans int
  keptCows int
  keptChicken int
}

var names [5]string = [5]string{"T", "D", "H", "C", "I"}

// Note that on a real die, the chances to get a death ray are twice the other type
// as there are two death ray pictures on each die.

func RollDie() int {
  val := rand.Intn(6)
  if val == 5 {
    val = deathray
  }
  return val
}

func (m *MartianState) Reset() {
  for i := range(m.dice) {
    m.dice[i].locked = false
  }
  m.keptTanks = 0
  m.keptDeathrays = 0
  m.keptChicken = 0
  m.keptCows = 0
  m.keptHumans = 0
}

func (m *MartianState) Roll() {
  for i := range(m.dice) {
    if !m.dice[i].locked {
      m.dice[i].value = RollDie();
    }
  }
}

// Debug print
func (m MartianState) PrintDice() {
  for i := range(m.dice) {
    // fmt.Println(i)
    fmt.Print(names[m.dice[i].value])
    if m.dice[i].locked {
      fmt.Print("*")
    } else {
      fmt.Print(" ")
    }
    fmt.Print(" ")
  }
  fmt.Println("\n")
}

func (m MartianState) PrintState() {
  fmt.Println("Tanks:", m.keptTanks)
  fmt.Print("Death Rays: ", m.keptDeathrays)
  if m.keptDeathrays < m.keptTanks {
    fmt.Println(" (WARNING)")
  } else {
    fmt.Println("")
  }
  fmt.Println("Humans:", m.keptHumans)
  fmt.Println("Cows:", m.keptCows)
  fmt.Println("chIcken:", m.keptChicken)
}

func (m *MartianState) NumDiceOfType(dieType int) int {
  count := 0
  for d := range(m.dice) {
    if m.dice[d].value == dieType {
      count++
    }
  }
  return count
}
func (m *MartianState) NumUnlockedDiceOfType(dieType int) int {
  count := 0
  for d := range(m.dice) {
    if m.dice[d].value == dieType && !m.dice[d].locked {
      count++
    }
  }
  return count
}

func (m *MartianState) LockDiceOfType(dieType int) {
  for d := range(m.dice) {
    if m.dice[d].value == dieType {
      m.dice[d].locked = true
    }
  }
}

func (m *MartianState) Keep(what int) (int, error) {
  if !m.CanKeepCreature(what) {
    return 0, errors.New("Can't keep this creature")
  }
  count := 0
  for d := range(m.dice) {
    if m.dice[d].value == what && !m.dice[d].locked {
      m.dice[d].locked = true
      count++
      switch (what) {
      case tank:
        m.keptTanks += 1
      case chicken:
        m.keptChicken += 1
      case cow:
        m.keptCows += 1
      case human:
        m.keptHumans += 1
      case deathray:
        m.keptDeathrays += 1
      }
    }
  }
  return count, nil
}

func PrintUsage() {
  fmt.Println("Did not understand.")
  fmt.Println("(Q to quit)")
}

func (m *MartianState) CanKeepCreature(value int) bool {
  switch value {
  case human:
    return m.NumUnlockedDiceOfType(value) > 0 && m.keptHumans == 0
  case chicken:
    return m.NumUnlockedDiceOfType(value) > 0 && m.keptChicken == 0
  case cow:
    return m.NumUnlockedDiceOfType(value) > 0 && m.keptCows == 0
  case deathray:
    return m.NumUnlockedDiceOfType(value) > 0;
  case tank:
    return true
  }
  return false
}

func (m *MartianState) ProcessCommand(cmd string) error {
  switch cmd {
  case "C":
    m.Keep(cow)
    return nil
  case "I":
    m.Keep(chicken)
    return nil
  case "H":
    m.Keep(human)
    return nil
  case "D":
    m.Keep(deathray)
    return nil
  }
  return errors.New("Bad command")
}

func (m *MartianState) CanMakeMove() bool {
  return m.CanKeepCreature(cow) || m.CanKeepCreature(chicken) || m.CanKeepCreature(human) || m.CanKeepCreature(deathray)
}

func (m *MartianState) ComputeScore() (int, int) {
  score := m.keptChicken + m.keptHumans + m.keptCows
  bonus := 0
  // Compute bonus points
  if m.keptChicken > 0 && m.keptHumans > 0 && m.keptCows > 0 {
    bonus = 3
  }
  return score, bonus
}

func main() {
  fmt.Println("Martian Dice")
  fmt.Println("============")
  reader := bufio.NewReader(os.Stdin)
  var m MartianState;
  curPlayer := 0

  rand.Seed(time.Now().UnixNano())

  numPlayers := 1
  fmt.Println("How many players?")
  _, err := fmt.Scanf("%d", &numPlayers)
  if err != nil {
    fmt.Println("??? Defaulting to 1 player.")
    numPlayers = 1
  }
  fmt.Println("Starting game with", numPlayers, "players.")
  fmt.Println("")
  p := make([]PlayerState, numPlayers)

  for {
    fmt.Println("\n==== Player", curPlayer + 1, "  Score:", p[curPlayer].Score, "=====")

    m.Reset()
    m.Roll()
    // Automatically lock tanks
    m.Keep(tank)
    for {
      m.PrintState()
      m.PrintDice()

      if !m.CanMakeMove() {
        fmt.Println("\nCan't make any more moves. End of round.\n")
        break
      }

      fmt.Println("Choose what to keep:")
      if m.CanKeepCreature(deathray) {
        fmt.Printf(" D (%d) ", m.NumUnlockedDiceOfType(deathray))
      }
      if m.CanKeepCreature(human) {
        fmt.Printf(" H (%d) ", m.NumUnlockedDiceOfType(human))
      }
      if m.CanKeepCreature(cow) {
        fmt.Printf(" C (%d) ", m.NumUnlockedDiceOfType(cow))
      }
      if m.CanKeepCreature(chicken) {
        fmt.Printf(" I (%d) ", m.NumUnlockedDiceOfType(chicken))
      }
      fmt.Println(" E (end turn)")
      command, _ := reader.ReadString('\n')
      command = strings.TrimSpace(strings.ToUpper(command))
      if command == "Q" {
        goto leave
      }
      err := m.ProcessCommand(command)
      if err == nil {
        // Valid move, let's roll dice.
        m.Roll()
        // Automatically lock tanks
        m.Keep(tank)
      }

      if command == "E" {
        break
      }
    }

    if m.keptTanks > m.keptDeathrays {
      fmt.Println("More tanks than death rays! You are dead, no points for you.")
    } else {
      score, bonus := m.ComputeScore()
      if bonus > 0 {
        fmt.Println("Got", bonus, "bonus points for rescuing all!")
      }
      p[curPlayer].Score += score + bonus
      fmt.Println("You scored", score+bonus, "points! Adding to your total, now", p[curPlayer].Score)
    }

    if p[curPlayer].Score >= 25 {
      fmt.Println("Player", curPlayer, "reached score 25 and won!")
      os.Exit(0)
    }

    curPlayer += 1
    curPlayer %= numPlayers
  }
  leave:
  ;
}
