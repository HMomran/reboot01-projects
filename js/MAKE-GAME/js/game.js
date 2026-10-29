import { GAME_CONFIG, GAME_ELEMENTS } from './config.js';
import { 
    updateLivesDisplay, 
    updateScoreDisplay, 
    updateRoundDisplay,
    updateTimerDisplay,
    showNotification,
    showGameOver,
    displayLeaderboard,
    hideLaunchMessage
} from './ui.js';
import { saveGameResult } from './storage.js';

export class Game {
    constructor() {
        this.reset();
    }

    reset() {
        this.currentPlayer = "";
        this.gameStarted = false;
        this.currentRound = 1;
        this.baseSpeed = GAME_CONFIG.BASE_SPEED;
        this.lives = GAME_CONFIG.STARTING_LIVES;
        this.score = 0;
        this.startTime = Date.now();
        this.elapsedTime = 0;
        this.timerStarted = false;
        this.oneMinuteWarning = false;
        this.isPaused = false;
        this.animationId = null;
        this.ballStuck = true;
        this.bricks = [];
        
        this.lastFrameTime = 0;
        this.deltaTime = 0;
        
        this.frameCount = 0;
        this.fpsTime = 0;
        this.currentFPS = 60;
        this.timerUpdateCounter = 0;
        
        this.ballX = GAME_CONFIG.BALL_START_X;
        this.ballY = GAME_CONFIG.BALL_START_Y;
        this.dx = GAME_CONFIG.BASE_SPEED;
        this.dy = -GAME_CONFIG.BASE_SPEED;
        
        this.paddleX = 0;
        this.leftPressed = false;
        this.rightPressed = false;
        
        this.borderWidth = GAME_CONFIG.GAME_WIDTH;
        this.borderHeight = GAME_CONFIG.GAME_HEIGHT;
        this.paddleWidth = GAME_CONFIG.PADDLE_WIDTH;
        this.paddleHeight = GAME_CONFIG.PADDLE_HEIGHT;

        // Dynamic layout (scaled with canvas height)
        this.paddleBottomOffset = 40;
        this.brickTopMarginPx = GAME_CONFIG.BRICK_TOP_MARGIN;
        this.minPlayGapPx = Math.round(GAME_CONFIG.GAME_HEIGHT * GAME_CONFIG.MIN_PLAY_GAP_RATIO);

        // Dynamic brick geometry (keeps visuals == collision)
        this.brickHeightPx = GAME_CONFIG.BRICK_HEIGHT;
        this.brickGapPx = GAME_CONFIG.BRICK_GAP;

    }

    start(playerName) {
        this.currentPlayer = playerName;
        this.gameStarted = true;

        this.updateDimensions();

        this.resetRound();
        this.createBricks();
        this.gameLoop();
    }

    handleResize() {
        if (!this.gameStarted) return;

        const oldBorderWidth = this.borderWidth;
        const oldBorderHeight = this.borderHeight;

        this.updateDimensions();

        const scaleX = oldBorderWidth ? (this.borderWidth / oldBorderWidth) : 1;
        const scaleY = oldBorderHeight ? (this.borderHeight / oldBorderHeight) : 1;

        // Keep paddle aligned after zoom/resize
        this.paddleX *= scaleX;
        this.paddleX = Math.max(0, Math.min(this.paddleX, this.borderWidth - this.paddleWidth));
        GAME_ELEMENTS.paddle.style.left = `${this.paddleX}px`;

        // Keep ball aligned after zoom/resize
        if (this.ballStuck) {
            this.ballX = this.paddleX + this.paddleWidth / 2 - GAME_CONFIG.BALL_SIZE / 2;
            this.ballY = this.borderHeight - this.paddleBottomOffset - this.paddleHeight - GAME_CONFIG.BALL_SIZE;
        } else {
            this.ballX *= scaleX;
            this.ballY *= scaleY;
        }

        // Re-sync brick collision geometry with current border size (and dynamic brick sizing)
        const availableWidth = this.borderWidth - (GAME_CONFIG.BRICK_LEFT_MARGIN * 2);
        const brickWidth = (availableWidth - (this.brickGapPx * (GAME_CONFIG.BRICK_COLS - 1))) / GAME_CONFIG.BRICK_COLS;

        for (let row = 0; row < GAME_CONFIG.BRICK_ROWS; row++) {
            for (let col = 0; col < GAME_CONFIG.BRICK_COLS; col++) {
                const idx = row * GAME_CONFIG.BRICK_COLS + col;
                const brick = this.bricks[idx];
                if (!brick) continue;

                brick.x = GAME_CONFIG.BRICK_LEFT_MARGIN + col * (brickWidth + this.brickGapPx);
                brick.y = this.brickTopMarginPx + row * (this.brickHeightPx + this.brickGapPx);
                brick.width = brickWidth;
                brick.height = this.brickHeightPx;

                // Keep DOM in sync with physics (alive bricks only)
                if (brick.el && brick.alive) {
                    brick.el.style.height = `${this.brickHeightPx}px`;
                }
            }
        }

        this.drawBall();
    }

    updateDimensions() {
        this.borderWidth = GAME_ELEMENTS.border.clientWidth;
        this.borderHeight = GAME_ELEMENTS.border.clientHeight;
        this.paddleWidth = GAME_ELEMENTS.paddle.offsetWidth;
        this.paddleHeight = GAME_ELEMENTS.paddle.offsetHeight;

        // Base offsets as ratios of current canvas height
        this.paddleBottomOffset = Math.round(this.borderHeight * GAME_CONFIG.PADDLE_BOTTOM_RATIO);
        this.brickTopMarginPx = Math.round(this.borderHeight * GAME_CONFIG.BRICK_TOP_RATIO);
        this.minPlayGapPx = Math.round(this.borderHeight * GAME_CONFIG.MIN_PLAY_GAP_RATIO);

        // Default brick geometry
        this.brickHeightPx = GAME_CONFIG.BRICK_HEIGHT;
        this.brickGapPx = GAME_CONFIG.BRICK_GAP;

        // Ensure we ALWAYS have enough space between brick-bottom and paddle-top.
        // If the canvas becomes too short (due to zoom/layout), shrink bricks first.
        const brickStackHeight0 =
            (GAME_CONFIG.BRICK_ROWS * this.brickHeightPx) +
            ((GAME_CONFIG.BRICK_ROWS - 1) * this.brickGapPx);

        const paddleTop0 = this.borderHeight - this.paddleBottomOffset - this.paddleHeight;
        const brickBottom0 = this.brickTopMarginPx + brickStackHeight0;
        let gap0 = paddleTop0 - brickBottom0;

        if (gap0 < this.minPlayGapPx) {
            // We need to free vertical space: shrink brick stack (height + gap) proportionally
            const requiredBrickBottomMax = paddleTop0 - this.minPlayGapPx;
            const maxStackHeight = Math.max(0, requiredBrickBottomMax - this.brickTopMarginPx);

            const scale = brickStackHeight0 > 0 ? Math.min(1, maxStackHeight / brickStackHeight0) : 1;

            // Clamp so bricks stay visible
            this.brickHeightPx = Math.max(10, Math.floor(GAME_CONFIG.BRICK_HEIGHT * scale));
            this.brickGapPx = Math.max(2, Math.floor(GAME_CONFIG.BRICK_GAP * scale));

            const brickStackHeight1 =
                (GAME_CONFIG.BRICK_ROWS * this.brickHeightPx) +
                ((GAME_CONFIG.BRICK_ROWS - 1) * this.brickGapPx);

            const paddleTop1 = this.borderHeight - this.paddleBottomOffset - this.paddleHeight;
            const brickBottom1 = this.brickTopMarginPx + brickStackHeight1;
            gap0 = paddleTop1 - brickBottom1;

            // If still too tight, push bricks up + paddle down (reduce offsets) while staying >= 0
            if (gap0 < this.minPlayGapPx) {
                let needed = this.minPlayGapPx - gap0;

                const reduceBrick = Math.min(this.brickTopMarginPx, Math.ceil(needed / 2));
                this.brickTopMarginPx -= reduceBrick;
                needed -= reduceBrick;

                const reducePaddle = Math.min(this.paddleBottomOffset, needed);
                this.paddleBottomOffset -= reducePaddle;
            }
        }

        // Apply DOM layout so visuals match physics
        GAME_ELEMENTS.paddle.style.bottom = `${this.paddleBottomOffset}px`;
        GAME_ELEMENTS.bricksEl.style.top = `${this.brickTopMarginPx}px`;
        GAME_ELEMENTS.bricksEl.style.gap = `${this.brickGapPx}px`;
    }

    resetRound() {
        this.ballX = this.borderWidth / 2 - GAME_CONFIG.BALL_SIZE / 2;
        this.ballY = this.borderHeight / 2 - GAME_CONFIG.BALL_SIZE / 2;
        this.dx = 0;
        this.dy = 0;
        this.ballStuck = true;
        
        this.paddleX = (this.borderWidth - this.paddleWidth) / 2;
        GAME_ELEMENTS.paddle.style.left = `${this.paddleX}px`;
        
        GAME_ELEMENTS.ball.style.left = `${this.ballX}px`;
        GAME_ELEMENTS.ball.style.top = `${this.ballY}px`;
        GAME_ELEMENTS.ball.style.display = "block";
        
        updateLivesDisplay(this.lives);
        updateScoreDisplay(this.score);
        updateRoundDisplay(this.currentRound, this.baseSpeed);
    }

    createBricks() {
        this.bricks = [];
        GAME_ELEMENTS.bricksEl.innerHTML = "";
        
        const availableWidth = this.borderWidth - (GAME_CONFIG.BRICK_LEFT_MARGIN * 2);
        const brickWidth = (availableWidth - (this.brickGapPx * (GAME_CONFIG.BRICK_COLS - 1))) / GAME_CONFIG.BRICK_COLS;

        for (let row = 0; row < GAME_CONFIG.BRICK_ROWS; row++) {
            for (let col = 0; col < GAME_CONFIG.BRICK_COLS; col++) {
                const brickEl = document.createElement("div");
                brickEl.className = "brick";
                brickEl.style.height = `${this.brickHeightPx}px`; // DOM matches physics
                GAME_ELEMENTS.bricksEl.appendChild(brickEl);

                const x = GAME_CONFIG.BRICK_LEFT_MARGIN + col * (brickWidth + this.brickGapPx);
                const y = this.brickTopMarginPx + row * (this.brickHeightPx + this.brickGapPx);

                this.bricks.push({
                    el: brickEl,
                    alive: true,
                    x: x,
                    y: y,
                    width: brickWidth,
                    height: this.brickHeightPx,
                    points: (GAME_CONFIG.BRICK_ROWS - row) * 10
                });
            }
        }
    }

    launchBall() {
        if (this.ballStuck && this.gameStarted && !this.isPaused) {
            this.ballStuck = false;
            hideLaunchMessage();
            this.dx = this.baseSpeed;
            this.dy = -this.baseSpeed;
            if (!this.timerStarted) {
                this.timerStarted = true;
                this.startTime = Date.now();
                this.elapsedTime = 0;
            }
        }
    }

    togglePause() {
        if (this.lives > 0 && this.gameStarted) {
            this.isPaused = !this.isPaused;
            return this.isPaused;
        }
        return false;
    }

    resume() {
        this.isPaused = false;
        this.startTime = Date.now() - this.elapsedTime;
    }

    showTimeWarning() {
        GAME_ELEMENTS.gameInfoNotificationEl.textContent = "⚠️ 1 MINUTE REMAINING!";
        GAME_ELEMENTS.gameInfoNotificationEl.hidden = false;
        setTimeout(() => {
            GAME_ELEMENTS.gameInfoNotificationEl.hidden = true;
        }, 3000);
    }

    updateTimer() {
        if (!this.isPaused && this.timerStarted) {
            this.elapsedTime = Date.now() - this.startTime;
            
            const timeLimit = GAME_CONFIG.TIME_LIMITS[this.currentRound - 1] || GAME_CONFIG.TIME_LIMITS[0];
            const remaining = timeLimit - this.elapsedTime;
            
            // Only end game if time is up AND there are still bricks remaining
            const activeBricks = this.bricks.filter(b => b.alive).length;
            if (remaining <= 0 && activeBricks > 0) {
                this.endGame(false);
                return;
            }
            
            if (remaining <= 60000 && !this.oneMinuteWarning) {
                this.oneMinuteWarning = true;
                this.showTimeWarning();
            }
        }
        this.timerUpdateCounter++;
        if (this.timerUpdateCounter >= 10) {
            const timeLimit = GAME_CONFIG.TIME_LIMITS[this.currentRound - 1] || GAME_CONFIG.TIME_LIMITS[0];
            const remaining = Math.max(0, timeLimit - this.elapsedTime);
            updateTimerDisplay(remaining);
            this.timerUpdateCounter = 0;
        }
    }

    updatePaddle() {
        const moveAmount = GAME_CONFIG.PADDLE_SPEED * this.deltaTime;
        
        if (this.leftPressed && this.paddleX > 0) {
            this.paddleX -= moveAmount;
        }

        if (this.rightPressed && this.paddleX < this.borderWidth - this.paddleWidth) {
            this.paddleX += moveAmount;
        }

        // Clamp paddle position to stay within boundaries
        this.paddleX = Math.max(0, Math.min(this.paddleX, this.borderWidth - this.paddleWidth));

        GAME_ELEMENTS.paddle.style.left = `${this.paddleX}px`;
    }

    updateBall() {
        if (this.ballStuck) {
            this.ballX = this.paddleX + this.paddleWidth / 2 - GAME_CONFIG.BALL_SIZE / 2;
            this.ballY = this.borderHeight - this.paddleBottomOffset - this.paddleHeight - GAME_CONFIG.BALL_SIZE;
        } else {
            const safeDelta = Math.min(this.deltaTime, 1.5);
            this.ballX += this.dx * safeDelta;
            this.ballY += this.dy * safeDelta;
        }
    }

    checkWallCollision() {
        if (this.ballX <= 0 || this.ballX + GAME_CONFIG.BALL_SIZE >= this.borderWidth) {
            this.dx = -this.dx;
            
            // Prevent ball from getting stuck in wall
            if (this.ballX <= 0) {
                this.ballX = 0;
            } else {
                this.ballX = this.borderWidth - GAME_CONFIG.BALL_SIZE;
            }
            
            const minHorizontalSpeed = 1.5;
            if (Math.abs(this.dx) < minHorizontalSpeed) {
                this.dx = this.dx > 0 ? minHorizontalSpeed : -minHorizontalSpeed;
            }
        }

        if (this.ballY <= 0) {
            this.dy = -this.dy;
            this.ballY = 0;
            
            // Ensure minimum horizontal speed when hitting ceiling
            const minHorizontalSpeed = 1.5;
            if (Math.abs(this.dx) < minHorizontalSpeed) {
                this.dx = this.dx >= 0 ? minHorizontalSpeed : -minHorizontalSpeed;
            }
        }
    }

    checkPaddleCollision() {
        const paddleTop = this.borderHeight - this.paddleBottomOffset - GAME_ELEMENTS.paddle.offsetHeight;
        const ballBottom = this.ballY + GAME_CONFIG.BALL_SIZE;
        //const ballCenterY = this.ballY + GAME_CONFIG.BALL_SIZE / 2;
        
        const safeDelta = Math.min(this.deltaTime, 1.5);
        const prevBallY = this.ballY - (this.dy * safeDelta);
        const prevBallBottom = prevBallY + GAME_CONFIG.BALL_SIZE;
        
        if (
            this.ballX + GAME_CONFIG.BALL_SIZE >= this.paddleX &&
            this.ballX <= this.paddleX + this.paddleWidth
        ) {
            if (
                prevBallBottom <= paddleTop &&
                ballBottom >= paddleTop &&
                this.dy > 0
            ) {
                this.dy = -this.dy;
                this.ballY = paddleTop - GAME_CONFIG.BALL_SIZE;
                
                const hitPos = (this.ballX + GAME_CONFIG.BALL_SIZE / 2) - (this.paddleX + this.paddleWidth / 2);
                const normalized = hitPos / (this.paddleWidth / 2);
                
                const currentSpeed = Math.sqrt(this.dx * this.dx + this.dy * this.dy);
                const maxAngle = 0.8;  
                const angle = normalized * maxAngle;
                
                this.dx = currentSpeed * Math.sin(angle);
                this.dy = -currentSpeed * Math.cos(angle);
                
                const maxSpeed = GAME_CONFIG.MAX_SPEED * this.currentRound;
                
                if (Math.abs(this.dx) > maxSpeed) {
                    this.dx = this.dx > 0 ? maxSpeed : -maxSpeed;
                }
                if (Math.abs(this.dy) > maxSpeed) {
                    this.dy = this.dy > 0 ? maxSpeed : -maxSpeed;
                }
                
                const minHorizontalSpeed = 1.5;
                if (Math.abs(this.dx) < minHorizontalSpeed) {
                    this.dx = this.dx >= 0 ? minHorizontalSpeed : -minHorizontalSpeed;
                }
                
                if (Math.abs(this.dy) < 2) {
                    this.dy = this.dy < 0 ? -2 : 2;
                }
            }
        }
    }

    checkBrickCollision() {
        const ballRight = this.ballX + GAME_CONFIG.BALL_SIZE;
        const ballBottom = this.ballY + GAME_CONFIG.BALL_SIZE;
        
        for (let i = 0; i < this.bricks.length; i++) {
            const brick = this.bricks[i];
            if (!brick.alive) continue;

            if (
                ballRight > brick.x &&
                this.ballX < brick.x + brick.width &&
                ballBottom > brick.y &&
                this.ballY < brick.y + brick.height
            ) {
                brick.alive = false;
                brick.el.style.visibility = "hidden";
                
                const ballCenterX = this.ballX + GAME_CONFIG.BALL_SIZE / 2;
                const ballCenterY = this.ballY + GAME_CONFIG.BALL_SIZE / 2;
                const brickCenterX = brick.x + brick.width / 2;
                const brickCenterY = brick.y + brick.height / 2;
                
                const diffX = Math.abs(ballCenterX - brickCenterX) / brick.width;
                const diffY = Math.abs(ballCenterY - brickCenterY) / brick.height;
                
                if (diffX > diffY) {
                    this.dx = -this.dx;
                    this.ballX = ballCenterX < brickCenterX ? brick.x - GAME_CONFIG.BALL_SIZE : brick.x + brick.width;
                } else {
                    this.dy = -this.dy;
                    this.ballY = ballCenterY < brickCenterY ? brick.y - GAME_CONFIG.BALL_SIZE : brick.y + brick.height;
                }
                
                // Ensure minimum horizontal speed after brick collision
                const minHorizontalSpeed = 1.5;
                if (Math.abs(this.dx) < minHorizontalSpeed) {
                    this.dx = this.dx >= 0 ? minHorizontalSpeed : -minHorizontalSpeed;
                }
                
                // Ensure minimum vertical speed
                if (Math.abs(this.dy) < 2) {
                    this.dy = this.dy >= 0 ? 2 : -2;
                }
                
                this.score += brick.points;
                updateScoreDisplay(this.score);
                
                this.checkWinCondition();
                break;
            }
        }
    }

    checkWinCondition() {
        const activeBricks = this.bricks.filter(b => b.alive).length;
        if (activeBricks === 0) {
            if (this.currentRound < GAME_CONFIG.MAX_ROUNDS) {
                this.currentRound++;
                this.baseSpeed = this.baseSpeed * 1.2;
                this.oneMinuteWarning = false;
                
                // Reset timer for new level
                this.elapsedTime = 0;
                this.startTime = Date.now();
                
                updateRoundDisplay(this.currentRound, this.baseSpeed);
                
                this.resetRound();
                this.createBricks();
                
                showNotification(`LEVEL ${this.currentRound}\\nSpeed Increased!`, GAME_CONFIG.ROUND_NOTIFICATION_DURATION);
            } else {
                this.endGame(true);
            }
        }
    }

    checkBallFall() {
        if (this.ballY > this.borderHeight) {
            this.lives--;
            updateLivesDisplay(this.lives);

            if (this.lives === 0) {
                this.endGame(false);
                return true;
            }

            this.paddleX = (this.borderWidth - this.paddleWidth) / 2;
            this.ballX = this.paddleX + this.paddleWidth / 2 - GAME_CONFIG.BALL_SIZE / 2;
            this.ballY = this.borderHeight - this.paddleBottomOffset - this.paddleHeight - GAME_CONFIG.BALL_SIZE;
            this.dx = 0;
            this.dy = 0;
            this.ballStuck = true;
            GAME_ELEMENTS.launchMessageEl.hidden = false;
            return false;
        }
        return false;
    }

    drawBall() {
        GAME_ELEMENTS.ball.style.left = `${this.ballX}px`;
        GAME_ELEMENTS.ball.style.top = `${this.ballY}px`;
        GAME_ELEMENTS.ball.style.display = "block";
        GAME_ELEMENTS.ball.style.visibility = "visible";
    }

    endGame(isVictory) {
        cancelAnimationFrame(this.animationId);
        saveGameResult(this.currentPlayer, this.score, this.elapsedTime);
        displayLeaderboard();
        showGameOver(this.currentPlayer, this.score, isVictory);
    }

    gameLoop(currentTime = 0) {
        if (this.lastFrameTime === 0) {
            this.lastFrameTime = currentTime;
            this.fpsTime = currentTime;
        }
        
        const frameTime = currentTime - this.lastFrameTime;
        this.lastFrameTime = currentTime;
        
        this.deltaTime = frameTime / 16.67;
        
        this.frameCount++;
        if (currentTime - this.fpsTime >= 2000) {
            this.currentFPS = Math.round(this.frameCount * 1000 / (currentTime - this.fpsTime));
            GAME_ELEMENTS.fpsEl.textContent = `FPS: ${this.currentFPS}`;
            this.frameCount = 0;
            this.fpsTime = currentTime;
        }
        
        if (this.isPaused) {
            this.lastFrameTime = 0;
            this.animationId = requestAnimationFrame((t) => this.gameLoop(t));
            return;
        }

        this.updateTimer();
        this.updatePaddle();
        this.updateBall();
        this.checkWallCollision();
        this.checkPaddleCollision();
        this.checkBrickCollision();
        
        if (this.checkBallFall()) {
            return;
        }

        this.drawBall();
        this.animationId = requestAnimationFrame((t) => this.gameLoop(t));
    }
}
