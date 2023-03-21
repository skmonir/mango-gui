import { Button, Card, Col, InputGroup, Row, Spinner } from "react-bootstrap";
import Form from "react-bootstrap/Form";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCode, faSave } from "@fortawesome/free-solid-svg-icons";
import { useEffect, useState } from "react";
import DataService from "../services/DataService.js";
import ViewCodeModal from "./modals/ViewCodeModal.jsx";
import ShowToast from "./Toast/ShowToast.jsx";
import Utils from "../Utils.js";

export default function Settings({ appState, setAppState }) {
  const [config, setConfig] = useState({
    workspaceDirectory: "",
    sourceDirectory: "",
    author: "",
    lang: "",
    compilationCommand: "",
    compilationArgs: "",
    templatePath: ""
  });

  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: ""
  });
  const [showToast, setShowToast] = useState(false);
  const [showCodeModal, setShowCodeModal] = useState(false);
  const [savingInProgress, setSavingInProgress] = useState(false);

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = () => {
    DataService.getConfig()
      .then(config => {
        saveConfigToUI(config);
      })
      .catch(e => {
        showToastMessage(
          "Error",
          "Oops! Something went wrong while fetching the config!"
        );
      });
  };

  const validate = () => {
    let errMessage = "";
    if (Utils.isStrNullOrEmpty(config.workspaceDirectory)) {
      errMessage += "Workspace directory path can't be empty\n";
    }
    if (Utils.isStrNullOrEmpty(config.compilationCommand)) {
      errMessage += "Compilation command can't be empty\n";
    }
    if (Utils.isStrNullOrEmpty(errMessage)) {
      return true;
    } else {
      showToastMessage("Error", errMessage);
      return false;
    }
  };

  const triggerSave = () => {
    if (validate()) {
      console.log("save triggerred");
      let configToSave = { ...appState.config };
      configToSave.author = config.author;
      configToSave.sourceDirectory = config.sourceDirectory;
      configToSave.workspaceDirectory = config.workspaceDirectory;
      configToSave.activeLanguage.lang = config.lang;
      configToSave.activeLanguage.compilationCommand =
        config.compilationCommand;
      configToSave.activeLanguage.compilationArgs = config.compilationArgs;
      configToSave.activeLanguage.templatePath = config.templatePath;
      let isFound = false;
      for (let i = 0; i < configToSave.languageConfigs.length; i++) {
        if (configToSave.languageConfigs[i].lang === config.lang) {
          isFound = true;
          configToSave.languageConfigs[i] = { ...configToSave.activeLanguage };
          break;
        }
      }
      if (!isFound) {
        configToSave.languageConfigs.push({ ...configToSave.activeLanguage });
      }
      console.log(configToSave);
      setSavingInProgress(true);
      DataService.updateConfig(configToSave)
        .then(config => {
          saveConfigToUI(config);
          showToastMessage("Success", "Settings saved successfully!");
        })
        .catch(e => {
          showToastMessage(
            "Error",
            "Oops! Something went wrong while saving the config!"
          );
        })
        .finally(() => setSavingInProgress(false));
    }
  };

  const saveConfigToUI = config => {
    console.log(config);
    setConfig({
      workspaceDirectory: config.workspaceDirectory,
      sourceDirectory: config.sourceDirectory,
      author: config.author,
      lang: config.activeLanguage.lang,
      compilationCommand: config.activeLanguage.compilationCommand,
      compilationArgs: config.activeLanguage.compilationArgs,
      templatePath: config.activeLanguage.templatePath
    });
    setAppState({ ...appState, config: config });
    console.log(appState.config);
  };

  const changeLanguage = lang => {
    setConfig({ ...config, lang: lang });
  };

  const showToastMessage = (variant, message) => {
    setShowToast(true);
    setToastMsgObj({
      variant: variant,
      message: message
    });
  };

  const checkDirectoryPathValidity = dirPath => {
    if (!Utils.isStrNullOrEmpty(dirPath)) {
      DataService.checkDirectoryPathValidity(window.btoa(dirPath)).then(
        resp => {
          if (resp.isExist === false) {
            showToastMessage("Error", `${dirPath} is not a valid directory`);
          }
        }
      );
    }
  };

  const checkFilePathValidity = filePath => {
    if (!Utils.isStrNullOrEmpty(filePath)) {
      DataService.checkFilePathValidity(window.btoa(filePath)).then(resp => {
        if (resp.isExist === false) {
          showToastMessage("Error", `${filePath} is not a valid file`);
        }
      });
    }
  };

  return (
    <div>
      <Card body bg="light">
        <Row>
          <Form.Group className="mb-3">
            <Form.Label>
              <strong>Workspace Directory</strong>
            </Form.Label>
            <Form.Control
              type="text"
              size="sm"
              autoCorrect="off"
              autoComplete="off"
              autoCapitalize="none"
              placeholder="Enter your workspace directory absolute path. All the testcases and sources will be saved here."
              value={config.workspaceDirectory}
              onChange={e =>
                setConfig({ ...config, workspaceDirectory: e.target.value })
              }
              onBlur={() =>
                checkDirectoryPathValidity(config.workspaceDirectory)
              }
            />
          </Form.Group>
        </Row>
        <Row>
          <Col sm={3}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>Language</strong>
              </Form.Label>
              <Form.Select
                size="sm"
                aria-label="Default select example"
                value={config.lang}
                onChange={e => changeLanguage(e.target.value)}
              >
                {/*<option value="">Select language</option>*/}
                <option value="c++">C++</option>
                {/*<option value="java">Java</option>*/}
              </Form.Select>
            </Form.Group>
          </Col>
          <Col sm={4}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>Compilation Command</strong>
              </Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="Example: g++"
                value={config.compilationCommand}
                onChange={e =>
                  setConfig({ ...config, compilationCommand: e.target.value })
                }
              />
            </Form.Group>
          </Col>
          <Col sm={5}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>Compilation Args</strong>
              </Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="Example: -std=c++20"
                value={config.compilationArgs}
                onChange={e =>
                  setConfig({ ...config, compilationArgs: e.target.value })
                }
              />
            </Form.Group>
          </Col>
        </Row>
        <Row>
          <Col sm={3}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>Author Name</strong>
              </Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="Enter your name"
                value={config.author}
                onChange={e => setConfig({ ...config, author: e.target.value })}
              />
            </Form.Group>
          </Col>
          <Col sm={9}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>Template File Path</strong>
              </Form.Label>
              <InputGroup className="mb-3">
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  value={config.templatePath}
                  onChange={e =>
                    setConfig({ ...config, templatePath: e.target.value })
                  }
                  onBlur={() => checkFilePathValidity(config.templatePath)}
                  placeholder={
                    "Template file ends with extension(.cpp). The template will be used to create source files."
                  }
                />
                <Button
                  size="sm"
                  variant="outline-success"
                  disabled={!config.templatePath}
                  onClick={() => setShowCodeModal(true)}
                >
                  <FontAwesomeIcon icon={faCode} /> View Code{" "}
                </Button>
              </InputGroup>
            </Form.Group>
          </Col>
        </Row>
        <Row>
          <Col md={{ span: 2, offset: 5 }}>
            <Row>
              <Col xs={12} className="d-flex justify-content-center">
                <Button
                  size="sm"
                  variant="outline-success"
                  onClick={() => triggerSave()}
                  disabled={savingInProgress}
                >
                  {!savingInProgress ? (
                    <FontAwesomeIcon icon={faSave} />
                  ) : (
                    <Spinner
                      as="span"
                      animation="grow"
                      size="sm"
                      role="status"
                      aria-hidden="true"
                    />
                  )}
                  {savingInProgress ? " Saving Settings" : " Save Settings"}
                </Button>
              </Col>
            </Row>
          </Col>
        </Row>
      </Card>
      {showToast && (
        <ShowToast toastMsgObj={toastMsgObj} setShowToast={setShowToast} />
      )}
      {showCodeModal && (
        <ViewCodeModal
          codePath={config.templatePath}
          setShowCodeModal={setShowCodeModal}
        />
      )}
    </div>
  );
}
