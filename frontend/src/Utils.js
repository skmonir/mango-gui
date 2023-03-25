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
}

export default new Utils();
