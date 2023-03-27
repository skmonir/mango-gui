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
  const [config, setConfig] = useState({});
  const [selectedLang, setSelectedLang] = useState("");
  const [selectedLangConfig, setSelectedLangConfig] = useState({});

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
        console.log(config);
        setConfig(config);
        setAppState({ ...appState, config: config });
        setSelectedLang(config?.activeLang);
        setSelectedLangConfig(config.langConfigs[config.activeLang]);
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
    if (Utils.isStrNullOrEmpty(errMessage)) {
      return true;
    } else {
      showToastMessage("Error", errMessage);
      return false;
    }
  };

  const triggerSave = async () => {
    console.log(config);
    if (validate()) {
      setSavingInProgress(true);
      DataService.updateConfig(updateLangConfigs())
        .then(config => {
          setAppState({ ...appState, config: config });
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

  const selectedLangChanged = lang => {
    if (!Utils.isStrNullOrEmpty(selectedLang)) {
      updateLangConfigs();
    }
    setSelectedLang(lang);
    setSelectedLangConfig(config.langConfigs[lang]);
  };

  const updateLangConfigs = () => {
    let updatedLangConfigs = {
      ...config.langConfigs
    };
    updatedLangConfigs[selectedLang] = { ...selectedLangConfig };
    const updatedConfig = {
      ...config,
      langConfigs: updatedLangConfigs
    };
    setConfig(updatedConfig);
    return updatedConfig;
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
          <Col xs={9}>
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
          </Col>
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
        </Row>
      </Card>
      <br />
      <Card body bg="light">
        <Row>
          <Col sm={2}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>Active Language</strong>
              </Form.Label>
              <Form.Select
                size="sm"
                aria-label="Default select example"
                value={config.activeLang}
                onChange={e =>
                  setConfig({ ...config, activeLang: e.currentTarget.value })
                }
              >
                <option value="cpp">CPP</option>
                <option value="java">Java</option>
                <option value="python">Python</option>
              </Form.Select>
            </Form.Group>
          </Col>
          <Col sm={2}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>Configure Language</strong>
              </Form.Label>
              <Form.Select
                size="sm"
                aria-label="Default select example"
                value={selectedLang}
                onChange={e => selectedLangChanged(e.currentTarget.value)}
              >
                <option value="cpp">CPP</option>
                <option value="java">Java</option>
                <option value="python">Python</option>
              </Form.Select>
            </Form.Group>
          </Col>
        </Row>
        {["java", "cpp"].includes(selectedLang) && (
          <Row>
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
                  value={selectedLangConfig.compilationCommand}
                  onChange={e =>
                    setSelectedLangConfig({
                      ...selectedLangConfig,
                      compilationCommand: e.target.value
                    })
                  }
                />
              </Form.Group>
            </Col>
            <Col sm={8}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Compilation Flags</strong>
                </Form.Label>
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  placeholder="Example: -std=c++20"
                  value={selectedLangConfig.compilationFlags}
                  onChange={e =>
                    setSelectedLangConfig({
                      ...selectedLangConfig,
                      compilationFlags: e.target.value
                    })
                  }
                />
              </Form.Group>
            </Col>
          </Row>
        )}
        {["java", "python"].includes(selectedLang) && (
          <Row>
            <Col sm={4}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Execution Command</strong>
                </Form.Label>
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  placeholder="Example: g++"
                  value={selectedLangConfig.executionCommand}
                  onChange={e =>
                    setSelectedLangConfig({
                      ...selectedLangConfig,
                      executionCommand: e.target.value
                    })
                  }
                />
              </Form.Group>
            </Col>
            <Col sm={8}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Execution Flags</strong>
                </Form.Label>
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  placeholder="Example: -std=c++20"
                  value={selectedLangConfig.executionFlags}
                  onChange={e =>
                    setSelectedLangConfig({
                      ...selectedLangConfig,
                      executionFlags: e.target.value
                    })
                  }
                />
              </Form.Group>
            </Col>
          </Row>
        )}
        <Row>
          <Col sm={12}>
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
                  value={selectedLangConfig.templatePath}
                  onChange={e =>
                    setSelectedLangConfig({
                      ...selectedLangConfig,
                      templatePath: e.target.value
                    })
                  }
                  onBlur={() =>
                    checkFilePathValidity(selectedLangConfig.templatePath)
                  }
                  placeholder={
                    "Template file ends with extension(.cpp/.java/.py). The sourceTemplateService will be used to create source files."
                  }
                />
                <Button
                  size="sm"
                  variant="outline-success"
                  disabled={!selectedLangConfig.templatePath}
                  onClick={() => setShowCodeModal(true)}
                >
                  <FontAwesomeIcon icon={faCode} /> View Code{" "}
                </Button>
              </InputGroup>
            </Form.Group>
          </Col>
        </Row>
      </Card>
      <br />
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
      {showToast && (
        <ShowToast toastMsgObj={toastMsgObj} setShowToast={setShowToast} />
      )}
      {showCodeModal && (
        <ViewCodeModal
          codePath={selectedLangConfig.templatePath}
          setShowCodeModal={setShowCodeModal}
        />
      )}
    </div>
  );
}
