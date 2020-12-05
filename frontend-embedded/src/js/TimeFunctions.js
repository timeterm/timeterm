.pragma library

Date.prototype.addHours = function (h) {
    this.setTime(this.getTime() + (h * 60 * 60 * 1000));
    return this;
}

Date.prototype.isFullHour = function () {
    return this.getMinutes() === 0 && this.getSeconds() === 0 && this.getMilliseconds() === 0;
}

Date.prototype.getMillisecondsInDay = function () {
    return this.getHours() * 3600 * 1000
        + this.getMinutes() * 60 * 1000
        + this.getSeconds() * 1000
        + this.getMilliseconds()
}

// startOfWeek calculates the start of the week, with Monday being the first day of the week.
Date.prototype.startOfWeek = function () {
    let newDate = new Date(this.getTime());

    const dayOfWeek = newDate.getDay();
    newDate.setDate(newDate.getDate() - (dayOfWeek === 0 ? 6 : dayOfWeek - 1));
    newDate.setHours(0, 0, 0, 0);

    return newDate;
}

Date.prototype.endOfWeek = function () {
    let date = this.startOfWeek();
    date.setDate(date.getDate() + 7);
    return date;
}
