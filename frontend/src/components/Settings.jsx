import {Alert, Button, Card, Col, InputGroup, Row, Spinner} from "react-bootstrap";
import Form from "react-bootstrap/Form";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCode, faSave} from "@fortawesome/free-solid-svg-icons";
import {useEffect, useState} from "react";
import DataService from "../services/DataService.js";
import ViewCodeModal from "./ViewCodeModal.jsx";

export default function Settings({appState, setAppState}) {
    const [config, setConfig] = useState({
        workspaceDirectory: '',
        sourceDirectory: '',
        author: '',
        lang: '',
        compilationCommand: '',
        compilationArgs: '',
        templatePath: ''
    });

    const [showCodeModal, setShowCodeModal] = useState(false);
    const [saveAlert, setSaveAlert] = useState('');
    const [savingInProgress, setSavingInProgress] = useState(false);

    useEffect(() => {
        fetchConfig()
    }, []);

    const fetchConfig = () => {
        DataService.getConfig().then(config => {
            saveConfigToUI(config);
        });
    };

    const triggerSave = () => {
        console.log("save triggerred");
        let configToSave = {...appState.config};
        configToSave.author = config.author;
        configToSave.sourceDirectory = config.sourceDirectory;
        configToSave.workspaceDirectory = config.workspaceDirectory;
        configToSave.activeLanguage.lang = config.lang;
        configToSave.activeLanguage.compilationCommand = config.compilationCommand;
        configToSave.activeLanguage.compilationArgs = config.compilationArgs;
        configToSave.activeLanguage.templatePath = config.templatePath;
        let isFound = false;
        for (let i = 0; i < configToSave.languageConfigs.length; i++) {
            if (configToSave.languageConfigs[i].lang === config.lang) {
                isFound = true;
                configToSave.languageConfigs[i] = {...configToSave.activeLanguage};
                break;
            }
        }
        if (!isFound) {
            configToSave.languageConfigs.push({...configToSave.activeLanguage});
        }
        console.log(configToSave);
        setSavingInProgress(true);
        DataService.updateConfig(configToSave).then(config => {
            saveConfigToUI(config);
            setSaveAlert('Settings saved successfully!');
        }).catch(e => {
            setSaveAlert('Oops! Something went wrong while saving the config!');
        }).finally(() => setSavingInProgress(false));
    };

    const saveConfigToUI = (config) => {
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
        setAppState({...appState, config: config});
        console.log(appState.config);
    }

    const changeLanguage = (lang) => {
        setConfig({...config, lang: lang});
    };

    return (
        <Card body bg="light">
            <Row>
                <Form.Group className="mb-3">
                    <Form.Label><strong>Workspace Directory</strong></Form.Label>
                    <Form.Control type="text" size="sm" placeholder="Enter your workspace directory absolute path"
                                  value={config.workspaceDirectory}
                                  onChange={(e) => setConfig({...config, workspaceDirectory: e.target.value})}/>
                </Form.Group>
            </Row>
            <Row>
                <Col sm={3}>
                    <Form.Group className="mb-3">
                        <Form.Label><strong>Language</strong></Form.Label>
                        <Form.Select size="sm" aria-label="Default select example" value={config.lang}
                                     onChange={(e) => changeLanguage(e.target.value)}>
                            {/*<option value="">Select language</option>*/}
                            <option value="c++">C++</option>
                            {/*<option value="java">Java</option>*/}
                        </Form.Select>
                    </Form.Group>
                </Col>
                <Col sm={4}>
                    <Form.Group className="mb-3">
                        <Form.Label><strong>Compilation Command</strong></Form.Label>
                        <Form.Control type="text" size="sm" placeholder="Example: g++"
                                      value={config.compilationCommand}
                                      onChange={(e) => setConfig({...config, compilationCommand: e.target.value})}/>
                    </Form.Group>
                </Col>
                <Col sm={5}>
                    <Form.Group className="mb-3">
                        <Form.Label><strong>Compilation Args</strong></Form.Label>
                        <Form.Control type="text" size="sm" placeholder="Example: -std=c++20"
                                      value={config.compilationArgs}
                                      onChange={(e) => setConfig({...config, compilationArgs: e.target.value})}/>
                    </Form.Group>
                </Col>
            </Row>
            <Row>
                <Col sm={3}>
                    <Form.Group className="mb-3">
                        <Form.Label><strong>Author Name</strong></Form.Label>
                        <Form.Control type="text" size="sm" placeholder="Enter your name" value={config.author}
                                      onChange={(e) => setConfig({...config, author: e.target.value})}/>
                    </Form.Group>
                </Col>
                <Col sm={9}>
                    <Form.Group className="mb-3">
                        <Form.Label><strong>Template File Path</strong></Form.Label>
                        <InputGroup className="mb-3">
                            <Form.Control type="text" size="sm" value={config.templatePath} onChange={(e) => setConfig({...config, templatePath: e.target.value})}/>
                            <Button size="sm" variant="outline-secondary" disabled={!config.templatePath} onClick={() => setShowCodeModal(true)}><FontAwesomeIcon icon={faCode}/> View Code </Button>
                            <Form.Text muted>Template file ends with extension(.cpp). The template will be used to create source files</Form.Text>
                        </InputGroup>
                    </Form.Group>
                </Col>
                {/*<Col sm={2}>*/}
                {/*</Col>*/}
            </Row>
            <Row>
                <Col md={{span: 4, offset: 5}}>
                    <Button size="sm" variant="outline-success" onClick={() => triggerSave()} disabled={savingInProgress}>
                        {
                            !savingInProgress ? (<FontAwesomeIcon icon={faSave}/>) : (<Spinner as="span" animation="grow" size="sm" role="status" aria-hidden="true"/>)
                        }
                        {savingInProgress ? ' Saving Settings' : ' Save Settings'}
                    </Button>
                </Col>
            </Row>
            {saveAlert !== '' && (
                <Row>
                    <Col>
                        <br/>
                        <Alert variant={saveAlert === 'Settings saved successfully!' ? 'success' : 'danger'} className="text-center">
                            {saveAlert}
                        </Alert>
                    </Col>
                </Row>
            )}
            {showCodeModal && (
                <ViewCodeModal codePath={config.templatePath} setShowCodeModal={setShowCodeModal}/>
            )}
        </Card>
    );
}