class Utils {
  isStrNullOrEmpty(str) {
    return (
      str === null || str === undefined || str.trim() === "" || str.length === 0
    );
  }

  isValidNum(n, min, max) {
    console.log("isValidNum: " + n);
    return (
      !isNaN(n) &&
      new RegExp("^[0-9]*$").test(n) &&
      min <= Number(n) &&
      Number(n) <= max
    );
  }

  dateStringToUiFormat(dateStr) {
    const formattedDate = new Date(dateStr).toLocaleString("en-US", {
      day: "numeric",
      month: "short",
      year: "numeric",
      hour: "numeric",
      minute: "2-digit"
    });
    return formattedDate;
  }

  dateToLocaleString(date) {
    return new Date(date).toLocaleString("en-US", {
      day: "numeric",
      month: "short",
      year: "numeric",
      hour: "numeric",
      minute: "2-digit"
    });
  }
}

export default new Utils();
