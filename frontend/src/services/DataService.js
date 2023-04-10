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

  addCustomProblem(addCustomProblemRequest) {
    return this.postData("problem/custom/add/", addCustomProblemRequest);
  }

  getConfig() {
    return this.getData("config/");
  }

  getEditorPreference() {
    return this.getData("config/editor");
  }

  updateEditorPreference(preference) {
    return this.putData("config/editor", preference);
  }

  updateConfig(config) {
    return this.putData("config/update", config);
  }

  resetConfig() {
    return this.getData("config/reset");
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

  generateSourceCode(path) {
    return this.getData("source/generate/" + path);
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
    return this.postData("testcase/random/input/generate/", req);
  }

  generateOutput(req) {
    return this.postData("testcase/random/output/generate/", req);
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

  checkDirectoryPathValidity(encodedPath) {
    return this.getData("misc/directory/check/" + encodedPath);
  }

  checkFilePathValidity(encodedPath) {
    return this.getData("misc/filepath/check/" + encodedPath);
  }

  openResource(openResourceRequest) {
    return this.putData("misc/resource/open/", openResourceRequest);
  }

  initApp() {
    return this.getData("misc/init");
  }

  getAppData() {
    return this.getData("appdata/");
  }

  scheduleParse(req) {
    return this.postData("schedule/parse/", req);
  }

  getParseScheduledTasks() {
    return this.getData("schedule/parse/");
  }

  removeParseScheduledTask(id) {
    return this.deleteData("schedule/parse/" + id);
  }
}

export default new DataService();
