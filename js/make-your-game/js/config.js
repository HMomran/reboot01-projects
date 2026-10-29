export const GAME_CONFIG = {
    GAME_WIDTH: 800,
    GAME_HEIGHT: 600,
    
    BALL_SIZE: 20,
    BALL_START_X: 390,
    BALL_START_Y: 300,
    
    PADDLE_WIDTH: 120,
    PADDLE_HEIGHT: 20,
    PADDLE_SPEED: 10,
    
    BRICK_ROWS: 5,
    BRICK_COLS: 10,
    BRICK_GAP: 8,
    BRICK_LEFT_MARGIN: 20,
    BRICK_TOP_MARGIN: 60,
    BRICK_HEIGHT: 25,

    PADDLE_BOTTOM_RATIO: 40 / 600,
    BRICK_TOP_RATIO: 60 / 600,
    MIN_PLAY_GAP_RATIO: 0.25,
    
    BASE_SPEED: 4,
    MAX_SPEED: 8,
    MAX_ROUNDS: 3,
    STARTING_LIVES: 3,
    TIME_LIMITS: [180000, 150000, 120000],
    
    NOTIFICATION_DURATION: 2000,
    ROUND_NOTIFICATION_DURATION: 2500
};

export const GAME_ELEMENTS = {
    paddle: document.getElementById("paddle"),
    border: document.getElementById("border"),
    gameInfo: document.getElementById("gameInfo"),
    ball: document.getElementById("ball"),
    livesEl: document.getElementById("lives"),
    roundInfoEl: document.getElementById("roundInfo"),
    scoreEl: document.getElementById("Score"),
    timeEl: document.getElementById("Time"),
    fpsEl: document.getElementById("fps"),
    startScreenEl: document.getElementById("startScreen"),
    usernameInput: document.getElementById("usernameInput"),
    startBtn: document.getElementById("startBtn"),
    usernameConfirmEl: document.getElementById("usernameConfirm"),
    confirmMessage: document.getElementById("confirmMessage"),
    confirmYes: document.getElementById("confirmYes"),
    confirmNo: document.getElementById("confirmNo"),
    pauseMenuEl: document.getElementById("pauseMenu"),
    resumeBtn: document.getElementById("resumeBtn"),
    restartBtn: document.getElementById("restartBtn"),
    gameOverMenuEl: document.getElementById("gameOverMenu"),
    playAgainBtn: document.getElementById("playAgainBtn"),
    exitBtn: document.getElementById("exitBtn"),
    finalScoreEl: document.getElementById("finalScore"),
    leaderboardBody: document.getElementById("leaderboardBody"),
    notificationEl: document.getElementById("notification"),
    notificationText: document.getElementById("notificationText"),
    launchMessageEl: document.getElementById("launchMessage"),
    bricksEl: document.getElementById("bricks"),
    gameOverTitle: document.getElementById("game-over-title"),
    gameInfoNotificationEl: document.getElementById("gameInfoNotification")
};
