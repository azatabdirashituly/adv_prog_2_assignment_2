package cmd

import "strings"

var keywords = []string{
    "soccer", "football", "match", "game", "pitch", "goal", "kick", "league",
    "tournament", "cup", "goalkeeper", "defender", "midfielder", "forward",
    "striker", "winger", "full-back", "centre-back", "world cup", "uefa",
    "champions league", "europa league", "premier league", "la liga", "serie a",
    "bundesliga", "ligue 1", "mls", "copa america", "euros", "pass", "shoot",
    "tackle", "dribble", "save", "header", "free kick", "penalty", "corner kick",
    "throw-in", "football boots", "jersey", "shorts", "shin guards", "goal nets",
    "football", "offside", "foul", "yellow card", "red card", "handball",
    "substitution", "var", "manchester united", "real madrid", "barcelona",
    "bayern munich", "juventus", "liverpool", "paris saint-germain", "chelsea",
    "arsenal", "ac milan", "messi", "ronaldo", "neymar", "mbappe", "salah",
    "modric", "lewandowski", "benzema", "haaland", "pogba", "formation",
    "strategy", "counter-attack", "set piece", "defense", "attack", "transition",
    "pressing", "high line", "tobyl",
}

func containsKeywords(request string) bool {
    words := strings.Fields(strings.ToLower(request))
    for _, word := range words {
        for _, keyword := range keywords {
            if strings.Contains(word, keyword) {
                return true
            }
        }
    }
    return false
}
