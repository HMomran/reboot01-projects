import { GAME_ELEMENTS, GAME_CONFIG } from './config.js';
import { getPlayers } from './storage.js';

export function showNotification(message, duration = GAME_CONFIG.NOTIFICATION_DURATION) {
    GAME_ELEMENTS.gameInfoNotificationEl.textContent = message;
    GAME_ELEMENTS.gameInfoNotificationEl.hidden = false;
    setTimeout(() => {
        GAME_ELEMENTS.gameInfoNotificationEl.hidden = true;
    }, duration);
}

export function updateLivesDisplay(lives) {
    GAME_ELEMENTS.livesEl.textContent = "Lives: " + lives;
}

export function updateScoreDisplay(score) {
    GAME_ELEMENTS.scoreEl.textContent = "SCORE: " + score;
}

export function updateRoundDisplay(round, speed) {
    GAME_ELEMENTS.roundInfoEl.textContent = `Round: ${round} | Speed: ${speed}x`;
}

export function updateTimerDisplay(elapsedTime) {
    const totalSeconds = Math.floor(elapsedTime / 1000);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    GAME_ELEMENTS.timeEl.textContent = `TIME: ${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
}

export function displayLeaderboard() {
    const players = getPlayers();
    
    players.sort((a, b) => {
        if (b.score !== a.score) {
            return b.score - a.score;
        }
        return a.time - b.time;
    });
    
    GAME_ELEMENTS.leaderboardBody.innerHTML = "";
    
    players.slice(0, 10).forEach((player, index) => {
        const row = document.createElement("tr");
        const minutes = Math.floor(player.time / 60);
        const seconds = player.time % 60;
        const timeStr = `${minutes}:${String(seconds).padStart(2, '0')}`;
        
        row.innerHTML = `
            <td>${index + 1}</td>
            <td>${player.name}</td>
            <td>${player.score}</td>
            <td>${timeStr}</td>
        `;
        GAME_ELEMENTS.leaderboardBody.appendChild(row);
    });
}

export function showUsernameConfirmation(name, onYes, onNo) {
    GAME_ELEMENTS.confirmMessage.textContent = `Found player "${name}". Are you this user?`;
    GAME_ELEMENTS.usernameConfirmEl.hidden = false;
    GAME_ELEMENTS.usernameInput.disabled = true;
    GAME_ELEMENTS.startBtn.disabled = true;
    
    const cleanup = () => {
        GAME_ELEMENTS.confirmYes.removeEventListener("click", handleYes);
        GAME_ELEMENTS.confirmNo.removeEventListener("click", handleNo);
        GAME_ELEMENTS.usernameConfirmEl.hidden = true;
        GAME_ELEMENTS.usernameInput.disabled = false;
        GAME_ELEMENTS.startBtn.disabled = false;
    };
    
    const handleYes = () => {
        cleanup();
        onYes();
    };
    
    const handleNo = () => {
        cleanup();
        onNo();
    };
    
    GAME_ELEMENTS.confirmYes.addEventListener("click", handleYes);
    GAME_ELEMENTS.confirmNo.addEventListener("click", handleNo);
}

export function showUsernameError() {
    GAME_ELEMENTS.usernameInput.style.border = "2px solid red";
    GAME_ELEMENTS.usernameInput.placeholder = "Name is required!";
    setTimeout(() => {
        GAME_ELEMENTS.usernameInput.style.border = "";
        GAME_ELEMENTS.usernameInput.placeholder = "Your Name";
    }, 2000);
}

export function showGameOver(playerName, score, isVictory = false) {
    document.body.classList.remove("game-active");
    GAME_ELEMENTS.finalScoreEl.textContent = isVictory 
        ? `${playerName}, You Win! Final Score: ${score}`
        : `${playerName}, Final Score: ${score}`;
    GAME_ELEMENTS.gameOverTitle.textContent = isVictory ? "VICTORY!" : "GAME OVER";
    GAME_ELEMENTS.gameOverMenuEl.hidden = false;
}

export function hidePauseMenu() {
    GAME_ELEMENTS.pauseMenuEl.hidden = true;
}

export function showPauseMenu() {
    GAME_ELEMENTS.pauseMenuEl.hidden = false;
}

export function hideStartScreen() {
    GAME_ELEMENTS.startScreenEl.hidden = true;
}

export function showGameBorder() {
    GAME_ELEMENTS.border.hidden = false;
    GAME_ELEMENTS.gameInfo.hidden = false;
}

export function hideGameBorder() {
    GAME_ELEMENTS.border.hidden = true;
    GAME_ELEMENTS.gameInfo.hidden = true;
}

export function hideGameOverMenu() {
    GAME_ELEMENTS.gameOverMenuEl.hidden = true;
}

export function showLaunchMessage() {
    GAME_ELEMENTS.launchMessageEl.hidden = false;
}

export function hideLaunchMessage() {
    GAME_ELEMENTS.launchMessageEl.hidden = true;
}
