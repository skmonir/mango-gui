import axios from "axios";
// import AuthService from './auth.service';

const BASE_URL = "http://localhost:3456/api/v1/";

class DataService {
  getData(url, params) {
    url = BASE_URL + url;

    let requestConfig = {
      params: params
    };

    return axios.get(url, requestConfig).then(response => response.data);
  }

  postData(url, payload) {
    url = BASE_URL + url;

    let requestConfig = {
      headers: {
        "Content-Type": "application/json"
      }
    };

    return axios
      .post(url, payload, requestConfig)
      .then(response => response.data);
  }

  putData(url, payload) {
    url = BASE_URL + url;

    let requestConfig = {
      headers: {
        "Content-Type": "application/json"
      }
    };

    return axios
      .put(url, payload, requestConfig)
      .then(response => response.data);
  }

  deleteData(url, data) {
    url = BASE_URL + url;

    let requestConfig = {
      headers: {
        "Content-Type": "application/json"
      },
      data: data
    };

    return axios.delete(url, requestConfig).then(response => response.data);
  }

  parse(encodedUrl) {
    return this.getData("parse/" + encodedUrl);
  }

  getProblem(path) {
    return this.getData("problem/" + path);
  }

  getProblemList(encodedUrl) {
    return this.getData("problem/" + encodedUrl);
  }

  getConfig() {
    return this.getData("config/");
  }

  updateConfig(config) {
    return this.putData("config/", config);
  }

  getCodeByMetadata(path) {
    return this.getData("code/" + path);
  }

  getCodeByPath(codeRequest) {
    return this.putData("code/", codeRequest);
  }

  updateCodeByFilePath(updateRequest) {
    return this.putData("code/update/", updateRequest);
  }

  updateCodeByProblemPath(prob_path, updateRequest) {
    return this.putData("code/update/" + prob_path, updateRequest);
  }

  openSourceByMetadata(path) {
    return this.getData("source/open/" + path);
  }

  openSourceByPath(openSourceRequest) {
    return this.putData("source/open/", openSourceRequest);
  }

  getTestcaseByFilePath(getTestcaseRequest) {
    return this.putData("testcase/custom/", getTestcaseRequest);
  }

  addCustomTest(addCustomTestRequest) {
    return this.postData("testcase/custom/add/", addCustomTestRequest);
  }

  updateCustomTest(updateCustomTestRequest) {
    return this.putData("testcase/custom/update/", updateCustomTestRequest);
  }

  deleteCustomTest(deleteCustomTestRequest) {
    return this.deleteData("testcase/custom/delete/", deleteCustomTestRequest);
  }

  generateRandomTests(req) {
    return this.postData("testcase/random/generate/", req);
  }

  runTest(path) {
    return this.getData("test/" + path);
  }

  getExecutionResult(path) {
    return this.getData("execresult/" + path);
  }

  getInputOutputDirectoriesByUrl(encodedUrl) {
    return this.getData("directories/" + encodedUrl);
  }
}

export default new DataService();
