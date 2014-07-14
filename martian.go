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
  "sort"
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

type DiceSlice []Die

func (dice DiceSlice) Swap(i, j int) {
  dice[i], dice[j] = dice[j], dice[i]
}

func (dice DiceSlice) Less(i, j int) bool {
  if (dice[i].locked && !dice[j].locked) {
    return true
  }
  if (!dice[i].locked && dice[j].locked) {
    return false
  }
  return dice[i].value < dice[j].value
}

func (dice DiceSlice) Len() int {
  return len(dice)
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

  // Automatically lock tanks
  m.Keep(tank)

  // Sort the dice.
  sort.Sort(DiceSlice(m.dice[:]))
}

func (m MartianState) PrintDice() {
  last := m.dice[0].value
  for i := range(m.dice) {
    if m.dice[i].value != last {
      last = m.dice[i].value
      fmt.Println()
    }
    if m.dice[i].locked {
      fmt.Printf("[%s]", names[m.dice[i].value])
    } else {
      fmt.Printf(" %s ", names[m.dice[i].value])
    }
    fmt.Print(" ")
  }
  fmt.Println("\n")
}

func (m MartianState) PrintState() {
  return
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
  fmt.Println("")
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

func (m *MartianState) Keep(what int) error {
  if !m.CanKeepCreature(what) {
    return errors.New("Can't keep/abduct this creature!")
  }
  for d := range(m.dice) {
    if m.dice[d].value == what && !m.dice[d].locked {
      m.dice[d].locked = true
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
  return nil
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
    return m.Keep(cow)
  case "I":
    return m.Keep(chicken)
  case "H":
    return m.Keep(human)
  case "D":
    return m.Keep(deathray)
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
  for x := 0; x < 10; x++ {
    fmt.Println("")
  }
  fmt.Println("Martian Dice")
  fmt.Println("============")
  fmt.Println("Abduct humans, cows and chicken! You get 1 point for each.")
  fmt.Println("Bonus points for abducting all three in one turn.")
  fmt.Println("You can only capture each kind once in a turn.")
  fmt.Println("You must end your turn with more death rays than enemy tanks.")
  fmt.Println("Type Q to quit at any time.")
  fmt.Println("Abbreviations: T = Enemy Tank, D = Death Ray, C = Cow, I = Chicken, H = Human.")
  fmt.Println("")

  reader := bufio.NewReader(os.Stdin)
  var m MartianState;
  curPlayer := 0

  rand.Seed(time.Now().UnixNano())

  numPlayers := 1
  fmt.Print("How many players?\n> ")
  _, err := fmt.Scanf("%d", &numPlayers)
  if err != nil {
    fmt.Println("??? Defaulting to 1 player.")
    numPlayers = 1
  }
  fmt.Println("Starting game with", numPlayers, "players.")
  fmt.Println("")
  p := make([]PlayerState, numPlayers)

  for {
    fmt.Printf("\n==== Player %d's turn (total score: %d) ====\n", curPlayer + 1, p[curPlayer].Score)

    m.Reset()
    m.Roll()
    for {
      m.PrintState()
      m.PrintDice()

      if !m.CanMakeMove() {
        fmt.Println("\nCan't make any more moves. End of round.\n")
        break
      }

      fmt.Print("Keep/abduct: ")
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
      fmt.Printf(" E (end turn)\n> ")
      command, _ := reader.ReadString('\n')
      fmt.Println()
      command = strings.TrimSpace(strings.ToUpper(command))
      if command == "Q" {
        goto leave
      }
      err := m.ProcessCommand(command)
      if err == nil {
        // Valid move, let's roll dice.
        m.Roll()
      } else {
        fmt.Println(err)
        fmt.Println()
      }

      if command == "E" {
        break
      }
    }

    if m.keptTanks > m.keptDeathrays {
      fmt.Printf("The %d tanks easily shot down your %d death rays! No points for you.\n\n", m.keptTanks, m.keptDeathrays)
    } else {
      score, bonus := m.ComputeScore()
      if score > 0 {
        fmt.Printf("You successfully abducted %d HUMANS, %d COWS and %d CHICKEN!\n", m.keptHumans, m.keptCows, m.keptChicken)
      } else {
        fmt.Printf("You failed to abduct any creatures! 0 points.")
      }
      if bonus > 0 {
        fmt.Println("Got", bonus, "bonus points for abducting all kinds!")
      }
      p[curPlayer].Score += score + bonus
      fmt.Println("You scored", score+bonus, "points! Adding to your total, now", p[curPlayer].Score)
      fmt.Println("")
    }

    fmt.Println("Press Enter to continue.")
    reader.ReadString('\n')

    if p[curPlayer].Score >= 25 {
      fmt.Println("Player", curPlayer + 1, "reached score 25 and won!")
      os.Exit(0)
    }

    curPlayer += 1
    curPlayer %= numPlayers
  }
  leave:
  ;
}
