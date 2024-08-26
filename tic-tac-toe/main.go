package main

import (
    "fmt"
)

const SIZE = 3

type GRID struct {
    grid [SIZE][SIZE]string
}

func PrintBoard(g *GRID, p string) {
    fmt.Printf("Current Player: %s\n", p)
    for i := 0; i < SIZE; i++ {
        for j := 0; j < SIZE; j++ {
            if g.grid[i][j] == "" {
                fmt.Print("- ")
            } else {
                fmt.Printf("%s ", g.grid[i][j])
            }
        }
        fmt.Println()
    }
}

func InsertSymbol(g *GRID, p string, row, col int) {
    row--
    col--
    g.grid[row][col] = p
}

func SwitchPlayer(currentPlayer string) string {
    if currentPlayer == "X" {
        return "O"
    }
    return "X"
}

func CheckSelection(g *GRID, row, col int) bool {
    if row > 3 || col > 3 || row < 1 || col < 1 {
        fmt.Println("Invalid Position. Please try again")
        return false
    }
    row--
    col--
    if g.grid[row][col] != "" {
        fmt.Println("Cell already selected")
        return false
    }
    return true
}

func CheckWin(g *GRID) bool {
    // Check rows
    for i := 0; i < SIZE; i++ {
        if g.grid[i][0] != "" && g.grid[i][0] == g.grid[i][1] && g.grid[i][1] == g.grid[i][2] {
            return true
        }
    }
    
    // Check columns
    for j := 0; j < SIZE; j++ {
        if g.grid[0][j] != "" && g.grid[0][j] == g.grid[1][j] && g.grid[1][j] == g.grid[2][j] {
            return true
        }
    }
    
    // Check diagonals
    if g.grid[0][0] != "" && g.grid[0][0] == g.grid[1][1] && g.grid[1][1] == g.grid[2][2] {
        return true
    }
    if g.grid[0][2] != "" && g.grid[0][2] == g.grid[1][1] && g.grid[1][1] == g.grid[2][0] {
        return true
    }
    
    return false
}

func main() {
    g := GRID{}
    currentPlayer := "X"
    var row, col int
    for {
        PrintBoard(&g, currentPlayer)
        fmt.Printf("Player %s's turn. Enter the row, col: ", currentPlayer)
        fmt.Scan(&row, &col)
        validSelection := CheckSelection(&g, row, col)
        if validSelection {
            InsertSymbol(&g, currentPlayer, row, col)
            if CheckWin(&g) {
                PrintBoard(&g, currentPlayer)
                fmt.Printf("Player %s wins!\n", currentPlayer)
                break
            }
            currentPlayer = SwitchPlayer(currentPlayer)
        }
    }
}