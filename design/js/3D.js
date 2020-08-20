function preload() {
    device = loadModel('assets/3d_files/device/device-without-RFID.obj');
    ophanging = loadModel('assets/3d_files/ophanging/ophanging.obj');
    loginScreen = loadImage('assets/images/login.svg');
    dayScreen = loadImage('assets/images/dag.svg');
    weekScreen = loadImage('assets/images/week.svg');
    z_urenScreen = loadImage('assets/images/z_uren.svg');
    card_in_hand = loadImage('assets/images/card_in_hand.svg');
    pointing_finger = loadImage('assets/images/pointing_finger.svg');
    wall = loadImage('assets/images/wall.jpg')
}

let fps = 60;
let capturer;
let btn;

function record() {
    capturer = new CCapture({
        format: "webm",
        framerate: fps
    });

    capturer.start();
    btn.textContent = "stop recording";
    btn.onclick = e => {
        capturer.stop();
        capturer.save();
        capturer = null;
        btn.textContent = "start recording";
        btn.onclick = record;
    };
}

function setup() {
    let sketch = createCanvas(1920, 1080, WEBGL); // 1280, 720 or 1920, 1080
    sketch.parent("container");

    timeOffset = 4000;
    moveTime = 1500;
    waitTime = 600;

    x = 0;
    y = 0;
    move = false;
    cursorView = true;
    screen = 0;
    milliSeconds = 0;

    if (!cursorView) {
        noCursor();
    }

    frameRate(fps);
    btn = document.getElementById("recordButton");
    btn.textContent = "start recording";
    btn.onclick = record;
}

var startMillis;

function draw() {
    milliSeconds = frameCount*1000/fps;

    background(255);

    walk();

    if (cursorView) {
        cursor();
    } else {
        noCursor();
    }

    drawDevice();
    drawWall();

    if (screen == 0) {
        drawLogin();
    } else if (screen == 1) {
        drawDay();
    } else if (screen == 2) {
        drawDay();
        drawZ_uren();
    } else if (screen == 3) {
        drawWeek();
    } else if (screen == 4) {
        drawWeek();
        drawZ_uren();
    }

    time = milliSeconds - timeOffset;

    if (time > 0) {
        if (time < moveTime) {
            drawCard(483 - time*300/moveTime , 0, 125 - time*75/moveTime);
        } else if (time < moveTime + waitTime) {
            drawCard(183, 0, 50);
        } else if (time < moveTime*2 + waitTime) {
            screen = 1;
            let timeMoved = time - (moveTime + waitTime);
            drawCard(183 + timeMoved*300/moveTime , 0, 50 + timeMoved*75/moveTime);
        }
    }

    drawHand();

    if (capturer) {
        capturer.capture(document.getElementById("defaultCanvas0"));
    }
}

function drawDevice() {
    pointLight(250, 250, 250, 0, 0, 500);
    noStroke();
    scale(1.8, -1.8);
    // let px = (width/2 - mouseX)/width;
    // let py = (height/2 - mouseY)/height;
    // if (move) {
    //     x = px;
    //     y = py;
    // }
    // rotateX(-y * TWO_PI);
    // rotateY(-x * TWO_PI);
    fill(200, 200, 200);
    
    push();
    translate(10.5, 0.5, -105);
    rotateZ(HALF_PI);
    model(ophanging);
    pop();

    translate(0, 10.5, 7.5);
    rotateX(-0.25*HALF_PI);
    model(device);
}

function drawWall() {
    noLights();

    push();
    rotateX(0.25*HALF_PI);
    translate(0, -10.5, -7.5);

    translate(0, 0, -105);
    scale(1, -1);
    texture(wall);
    plane(3300, 1800);

    push();
    translate(1650, 0, 1550);
    rotateY(HALF_PI);
    texture(wall);
    plane(3500, 1800);
    pop();

    translate(0, 900, 1650);
    rotateX(HALF_PI);
    fill(200, 200, 220);
    plane(3300, 3300);

    translate(0, 0, 1800);
    fill(255, 255, 230);
    plane(3300, 3300);

    pop();
}

function walk() {
    timeToStart = 1000;
    time = milliSeconds - timeToStart;
    startX = -3000;
    startY = -750;
    startZ = 2200;
    endX = 15;
    endY = -179-1;
    endZ = 500+47.5;
    startCenterX = 0;
    startCenterY = -750;
    startCenterZ = 0;
    endCenterX = 15;
    endCenterY = -40;
    endCenterZ = 47.5;
    timeToMove = 3000;

    if (time < 0) {
        camera(startX, startY, startZ, 0, startCenterY, 0, 0, 1, 0);
    } else if (time < timeToMove) {
        currX = lerp(startX, endX, time/timeToMove);
        currY = lerp(startY, endY, time/timeToMove);
        currZ = lerp(startZ, endZ, time/timeToMove);
        currCenterX = lerp(startCenterX, endCenterX, time/timeToMove);
        currCenterY = lerp(startCenterY, endCenterY, time/timeToMove);
        currCenterZ = lerp(startCenterZ, endCenterZ, time/timeToMove);
        camera(currX, currY, currZ, currCenterX, currCenterY, currCenterZ, 0, 1, 0);
    } else {
        camera(endX, endY, endZ, endCenterX, endCenterY, endCenterZ, 0, 1, 0);
    }
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

function drawWeek() {
    noLights();

    push();
    translate(-11, -1, 47.5);
    scale(1, -1);
    texture(weekScreen);
    plane(175*weekScreen.width/weekScreen.height, 175);
    pop();
}

function drawZ_uren() {
    noLights();

    push();
    translate(-11, -1, 47.5);
    scale(1, -1);
    fill('rgba(0, 0, 0, 0.25)')
    plane(175*weekScreen.width/weekScreen.height, 175);
    pop();

    push();
    translate(16.34, -7.84, 47.5);
    scale(1, -1);
    texture(z_urenScreen);
    plane(138.5*z_urenScreen.width/z_urenScreen.height, 138.5);
    pop();
}

function drawCard(x, y, z) {
    noLights();

    push();
    translate(x+40, y, z);
    scale(1, -1);
    texture(card_in_hand);
    plane(180, 180);
    pop();
}

function drawHand() {
    noLights();
    push();
    scale(1, -1);

    startX = 400;
    startY = 400;
    startZ = 150;

    endX1 = 25; // Week button
    endY1 = 155;
    endZ1 = 50;
    endX2 = 160; // Z_uren button in week
    endY2 = 135;
    endZ2 = 50;
    endX3 = 270; // Select z_uur
    endY3 = 130;
    endZ3 = 50;
    endX4 = 25; // Log-out button
    endY4 = 230;
    endZ4 = 50;
    duration = 1000;

    timeNow = milliSeconds - 8000;


    if (timeNow > 0 && timeNow <= duration) { // From Day to Z_uren
        percentage = (timeNow) / duration;

        currX = lerp(startX, endX2, percentage);
        currY = lerp(startY, endY2, percentage);
        currZ = lerp(startZ, endZ2, percentage);

        translate(currX, currY, currZ);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration && timeNow <= duration + 500) {
        translate(endX2, endY2, endZ2);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration + 500 && timeNow <= duration*2 + 500) {
        screen = 4;
        percentage = (timeNow-duration-500) / duration;

        currX = lerp(endX2, startX, percentage);
        currY = lerp(endY2, startY, percentage);
        currZ = lerp(endZ2, startZ, percentage);

        translate(currX, currY, currZ);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*2 + 2000 && timeNow <= duration*3 + 2000) { // From Z_uren
        percentage = (timeNow-duration*2-2000) / duration;

        currX = lerp(startX, endX3, percentage);
        currY = lerp(startY, endY3, percentage);
        currZ = lerp(startZ, endZ3, percentage);

        translate(currX, currY, currZ);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*3 + 2000 && timeNow <= duration*3 + 2500) {
        translate(endX3, endY3, endZ3);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*3 + 2500 && timeNow <= duration*4 + 2500) {
        screen = 1;
        percentage = (timeNow-duration*3-2500) / duration;

        currX = lerp(endX3, startX, percentage);
        currY = lerp(endY3, startY, percentage);
        currZ = lerp(endZ3, startZ, percentage);

        translate(currX, currY, currZ);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*4 + 4000 && timeNow <= duration*5 + 4000) { // From Day to Week
        percentage = (timeNow-duration*4-4000) / duration;

        currX = lerp(startX, endX1, percentage);
        currY = lerp(startY, endY1, percentage);
        currZ = lerp(startZ, endZ1, percentage);

        translate(currX, currY, currZ);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*5 + 4000 && timeNow <= duration*5 + 4500) {
        translate(endX1, endY1, endZ1);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*5 + 4500 && timeNow <= duration*6 + 4500) {
        screen = 3;
        percentage = (timeNow-duration*5-4500) / duration;

        currX = lerp(endX1, startX, percentage);
        currY = lerp(endY1, startY, percentage);
        currZ = lerp(endZ1, startZ, percentage);

        translate(currX, currY, currZ);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*6 + 6000 && timeNow <= duration*7 + 6000) { // From Week to Z_uren
        percentage = (timeNow-duration*6-6000) / duration;

        currX = lerp(startX, endX2, percentage);
        currY = lerp(startY, endY2, percentage);
        currZ = lerp(startZ, endZ2, percentage);

        translate(currX, currY, currZ);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*7 + 6000 && timeNow <= duration*7 + 6500) {
        translate(endX2, endY2, endZ2);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*7 + 6500 && timeNow <= duration*8 + 6500) {
        screen = 4;
        percentage = (timeNow-duration*7-6500) / duration;

        currX = lerp(endX2, startX, percentage);
        currY = lerp(endY2, startY, percentage);
        currZ = lerp(endZ2, startZ, percentage);

        translate(currX, currY, currZ);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*8 + 8000 && timeNow <= duration*9 + 8000) { // From Z_uren
        percentage = (timeNow-duration*8-8000) / duration;

        currX = lerp(startX, endX3, percentage);
        currY = lerp(startY, endY3, percentage);
        currZ = lerp(startZ, endZ3, percentage);

        translate(currX, currY, currZ);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*9 + 8000 && timeNow <= duration*9 + 8500) {
        translate(endX3, endY3, endZ3);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*9 + 8500 && timeNow <= duration*10 + 8500) {
        screen = 3;
        percentage = (timeNow-duration*9-8500) / duration;

        currX = lerp(endX3, startX, percentage);
        currY = lerp(endY3, startY, percentage);
        currZ = lerp(endZ3, startZ, percentage);

        translate(currX, currY, currZ);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*10 + 10000 && timeNow <= duration*11 + 10000) { // To Log-out
        percentage = (timeNow-duration*10-10000) / duration;

        currX = lerp(startX, endX4, percentage);
        currY = lerp(startY, endY4, percentage);
        currZ = lerp(startZ, endZ4, percentage);

        translate(currX, currY, currZ);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*11 + 10000 && timeNow <= duration*11 + 10500) {
        translate(endX4, endY4, endZ4);
        texture(pointing_finger);
        plane(340, 340);
    } else if (timeNow > duration*11 + 10500 && timeNow <= duration*12 + 10500) {
        screen = 0;
        percentage = (timeNow-duration*11-10500) / duration;

        currX = lerp(endX4, startX, percentage);
        currY = lerp(endY4, startY, percentage);
        currZ = lerp(endZ4, startZ, percentage);

        translate(currX, currY, currZ);
        texture(pointing_finger);
        plane(340, 340);
    }

    pop();
}

function keyPressed() {
    if (key == 'f' || key == 'F') {
        move = !move;
    } else if (key == 'c' || key == 'C') {
        cursorView = !cursorView;
    } else if (key == 'r' || key == 'R') {
        if(capturer) {
            capturer.stop();
            capturer.save();
            capturer = null;
            btn.textContent = "start recording";
            btn.onclick = record;
        } else {
            record();
        }
    }
}