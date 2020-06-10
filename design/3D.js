function preload() {
    device = loadModel('assets/device/device-without-RFID.obj');
    loginScreen = loadImage('assets/images/login.svg');
    dayScreen = loadImage('assets/images/day.svg');
    pas = loadImage('assets/images/pas.svg');
}

function setup() {
    let sketch = createCanvas(windowWidth, windowHeight-150, WEBGL);
    sketch.parent("container");

    timeOffset = 2000;
    moveTime = 1500;
    waitTime = 600;

    x = 0;
    y = 0;
    move = false;
    cursorView = false;
    signedIn = false;

    if (!cursorView) {
        noCursor();
    }
}

function draw() {
    background(255);

    if (cursorView) {
        cursor();
    } else {
        noCursor();
    }

    drawDevice();

    if (!signedIn) {
        drawLogin();
    } else {
        drawDay();
    }

    time = millis() - timeOffset;

    if (time > 0) {
        if (time < moveTime) {
            drawCard(483 - time*300/moveTime , 0, 125 - time*75/moveTime);
        } else if (time < moveTime + waitTime) {
            drawCard(183, 0, 50);
        } else if (time < moveTime*2 + waitTime) {
            signedIn = true;
            let timeMoved = time - (moveTime + waitTime);
            drawCard(183 + timeMoved*300/moveTime , 0, 50 + timeMoved*75/moveTime);
        }
    }
}

function drawDevice() {
    pointLight(250, 250, 250, 0, 0, 500);
    noStroke();
    scale(1.8, -1.8);
    let px = (width/2 - mouseX)/width;
    let py = (height/2 - mouseY)/height;
    if (move) {
        x = px;
        y = py;
    }
    rotateX(-y * TWO_PI);
    rotateY(-x * TWO_PI);
    fill(200, 200, 200);
    model(device);
}

function drawLogin() {
    noLights();

    push();
    translate(-11, -1, 47.5);
    scale(1, -1);
    texture(loginScreen);
    plane(175*loginScreen.width/loginScreen.height, 175);
    pop();
}

function drawDay() {
    noLights();

    push();
    translate(-11, -1, 47.5);
    scale(1, -1);
    texture(dayScreen);
    plane(175*dayScreen.width/dayScreen.height, 175);
    pop();
}

function drawCard(x, y, z) {
    noLights();

    push();
    translate(x, y, z);
    scale(1, -1);
    texture(pas);
    plane(85.60, 53.98);
    pop();
}

function keyPressed() {
    if (key == 'f' || key == 'F') {
        move = !move;
    } else if (key == 'c' || key == 'C') {
        cursorView = !cursorView;
    } else {
        timeOffset = millis();
        signedIn = false;
    }
}