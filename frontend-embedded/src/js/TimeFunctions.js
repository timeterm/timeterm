.pragma library

Date.prototype.addHours = function (h) {
    this.setTime(this.getTime() + (h * 60 * 60 * 1000));
    return this;
}

Date.prototype.addDays = function (d) {
    this.setTime(this.getTime() + (d * 24 * 60 * 60 * 1000));
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

Date.prototype.getWeek = function() {
  var date = new Date(this.getTime());
  date.setHours(0, 0, 0, 0);
  // Thursday in current week decides the year.
  date.setDate(date.getDate() + 3 - (date.getDay() + 6) % 7);
  // January 4 is always in week 1.
  var week1 = new Date(date.getFullYear(), 0, 4);
  // Adjust to Thursday in week 1 and count number of weeks from date to week1.
  return 1 + Math.round(((date.getTime() - week1.getTime()) / 86400000
                        - 3 + (week1.getDay() + 6) % 7) / 7);
}
