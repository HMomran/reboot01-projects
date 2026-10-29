import { GAME_ELEMENTS } from './config.js';
import { Game } from './game.js';
import { checkUsername } from './storage.js';
import { 
    showUsernameConfirmation, 
    showUsernameError, 
    hideStartScreen, 
    showGameBorder,
    hidePauseMenu,
    showPauseMenu,
    hideGameOverMenu,
    showLaunchMessage
} from './ui.js';

const game = new Game();

window.addEventListener("resize", () => {
    game.handleResize();
});

async function handleUsernameCheck(username) {
    const { existing, players } = checkUsername(username);
    
    return new Promise((resolve) => {
        if (existing) {
            const handleYes = () => {
                resolve(username);
            };
            
            const handleNo = () => {
                let counter = 1;
                let newName = `${username}${counter}`;
                while (players.find(p => p.name === newName)) {
                    counter++;
                    newName = `${username}${counter}`;
                }
                resolve(newName);
            };
            
            showUsernameConfirmation(username, handleYes, handleNo);
        } else {
            resolve(username);
        }
    });
}

GAME_ELEMENTS.startBtn.addEventListener("click", async () => {
    const username = GAME_ELEMENTS.usernameInput.value.trim();
    
    if (!username) {
        showUsernameError();
        return;
    }
    
    const currentPlayer = await handleUsernameCheck(username);
    
    hideStartScreen();
    document.body.classList.add("game-active");
    showGameBorder();
    showLaunchMessage();
    
    game.reset();
    game.start(currentPlayer);
});

document.addEventListener("keydown", (e) => {
    if (e.key === "Escape") {
        e.preventDefault();
        const paused = game.togglePause();
        if (paused) {
            showPauseMenu();
        } else {
            hidePauseMenu();
            game.resume();
        }
    }
    
    if (e.key === "ArrowLeft") game.leftPressed = true;
    if (e.key === "ArrowRight") game.rightPressed = true;
    
    if (e.key === "ArrowUp") {
        game.launchBall();
    }
});

document.addEventListener("keyup", (e) => {
    if (e.key === "ArrowLeft") game.leftPressed = false;
    if (e.key === "ArrowRight") game.rightPressed = false;
});

GAME_ELEMENTS.resumeBtn.addEventListener("click", () => {
    game.isPaused = false;
    hidePauseMenu();
    game.resume();
});

GAME_ELEMENTS.restartBtn.addEventListener("click", () => {
    if (game.animationId) {
        cancelAnimationFrame(game.animationId);
    }
    
    game.isPaused = false;
    hidePauseMenu();
    document.body.classList.add("game-active");
    showGameBorder();
    showLaunchMessage();
    
    const playerName = game.currentPlayer;
    game.reset();
    game.start(playerName);
});

GAME_ELEMENTS.playAgainBtn.addEventListener("click", () => {
    if (game.animationId) {
        cancelAnimationFrame(game.animationId);
    }
    
    hideGameOverMenu();
    document.body.classList.add("game-active");
    showGameBorder();
    showLaunchMessage();
    
    const playerName = game.currentPlayer;
    game.reset();
    game.start(playerName);
});

GAME_ELEMENTS.exitBtn.addEventListener("click", () => {
    location.reload();
});
