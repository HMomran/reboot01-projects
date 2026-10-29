const STORAGE_KEY = "brickBreakerPlayers";

export function getPlayers() {
    const data = localStorage.getItem(STORAGE_KEY);
    return data ? JSON.parse(data) : [];
}

export function savePlayers(players) {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(players));
}

export function saveGameResult(playerName, score, elapsedTime) {
    const players = getPlayers();
    const totalSeconds = Math.floor(elapsedTime / 1000);
    
    players.push({
        name: playerName,
        score: score,
        time: totalSeconds,
        date: new Date().toISOString()
    });
    
    savePlayers(players);
}

export function checkUsername(name) {
    const players = getPlayers();
    const existing = players.find(p => p.name === name);
    
    return { existing, players };
}
