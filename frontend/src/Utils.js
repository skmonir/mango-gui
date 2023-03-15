class Utils {
  isStrNullOrEmpty(str) {
    return (
      str === null || str === undefined || str.trim() === "" || str.length === 0
    );
  }
}

export default new Utils();
