.pragma library

Date.prototype.addHours = function(h) {
    this.setTime(this.getTime() + (h*60*60*1000));
    return this;
}

Date.prototype.isFullHour = function() {
    if (this.getMinutes() === 0 && this.getSeconds() === 0 && this.getMilliseconds() === 0) {
        return true;
    }
    return false;
}
