import {Accordion, Alert, Button, Card, Col, Row, Spinner, Table} from "react-bootstrap";
import Form from "react-bootstrap/Form";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCode, faFileCode, faTasks, faTerminal} from "@fortawesome/free-solid-svg-icons";
import {useEffect, useState} from "react";
import SocketClient from "../socket/SocketClient.js";
import DataService from "../services/DataService.js";
import ViewCodeModal from "./ViewCodeModal.jsx";

export default function Tester({appState, setAppState}) {
    const socketClient = new SocketClient();

    const [initAlert, setInitAlert] = useState(false);
    const [loadingInProgress, setLoadingInProgress] = useState(false);
    const [problemList, setProblemList] = useState([]);
    const [selectedProblem, setSelectedProblem] = useState({});
    const [selectedProblemMetadata, setSelectedProblemMetadata] = useState('');
    const [showCodeModal, setShowCodeModal] = useState(false);
    const [currentCodePath, setCurrentCodePath] = useState('');


    useEffect(() => {
        let socketConn = socketClient.initSocketConnection("test_problem_event", updateTestStatusFromSocket);
        return () => {
            socketConn.close()
        }
    }, []);

    const getProblemList = () => {
        setLoadingInProgress(true);
        setTimeout(() => {
            DataService.getProblemList(window.btoa(appState.url)).then(data => {
                setLoadingInProgress(false);
                setProblemList(data ? data : []);
                if (data && data.length > 0) {
                    getSelectedProblemInfo(data[0].platform + '/' + data[0].contestId + '/' + data[0].label);
                } else {
                    getSelectedProblemInfo('');
                }
            });
        }, 0);
    };

    const getSelectedProblemInfo = (metadata) => {
        setSelectedProblemMetadata(metadata);
        if (metadata && metadata.length > 0) {
            // setLoadingInProgress(true);
            DataService.getExecutionResult(metadata).then(data => {
                setSelectedProblem(data);
                setLoadingInProgress(false);
                // setInitAlert(true);
            });
        } else {
            setSelectedProblem({});
        }
    };

    const openSource = () => {
        DataService.openSource({filePath: currentCodePath}).then(resp => {
            console.log(resp);
        })
    };

    const runTest = () => {
        DataService.testCode(selectedProblem).then(data => {
            console.log(data);
            // setCurrentProblem(data);
        });
    };

    const prepareCurrentCodePath = (path) => {
        const values = path.split('/');
        let codepath = appState.config.workspaceDirectory.trimEnd() + '/' + values[0] + '/' + values[1] + '/source/' + values[2] + '.cpp';
        codepath = codepath.replaceAll('//', '/');
        setCurrentCodePath(codepath);
    }

    const updateTestStatusFromSocket = (data) => {
        // console.log(data);
        if (data) {
            // setCurrentProblem(data);
        }
    }

    const getTestStatusText = (testStatus) => {
        const errorList = ['Source fileService not found!', 'Compilation error!', 'Binary fileService not found!'];

        if (errorList.includes(testStatus)) {
            return <strong style={{color: 'red'}}>{testStatus}</strong>
        } else if (testStatus === 'Compilation Successful!') {
            return <strong style={{color: 'green'}}>{testStatus}</strong>
        } else if (testStatus.endsWith('Tests Passed')) {
            return <strong style={{color: currentProblem.isPassed ? 'green' : 'red'}}>{testStatus}</strong>
        } else {
            return <strong style={{color: 'purple'}}>{testStatus}</strong>
        }
    };

    const getVerdict = (execRes) => {
        if (execRes.message === 'running') {
            return <Spinner animation="border" variant="primary" size="sm"/>;
        } else if (execRes.verdict) {
            if (execRes.verdict === 'OK') {
                return <pre style={{color: 'green'}}>
                    <img src={`/src/assets/icons/${execRes.verdict}.svg`} style={{maxWidth: '30px'}}/> <strong>{execRes.verdict}</strong>
                </pre>
            } else {
                return <pre style={{color: 'red'}}>
                    <img src={`/src/assets/icons/${execRes.verdict}.svg`} style={{maxWidth: '30px'}}/> <strong>{execRes.verdict}</strong>
                </pre>
            }
        }
    }

    const disableActionButtons = () => {
        return !appState.url || appState.url === '' || loadingInProgress || !appState.config.workspaceDirectory;
    }

    // const getExecutionTable = () => {
    //     if (currentProblem && currentProblem.testcases && currentProblem.testcases.length > 0) {
    //         return (
    //             <div>
    //                 <div style={{height: '72vh', overflowY: "auto", overflowX: "auto"}}>
    //                     <Table bordered responsive="sm" size="sm">
    //                         <thead style={{position: 'sticky', top: 0, zIndex: 1, background: '#fff'}}>
    //                         <tr className="text-center">
    //                             <th>INPUT</th>
    //                             <th>OUTPUT</th>
    //                             <th>EXPECTED</th>
    //                             <th>VERDICT</th>
    //                             <th>TIME</th>
    //                             <th>MEMORY</th>
    //                         </tr>
    //                         </thead>
    //                         <tbody>
    //                         {
    //                             currentProblem.testcases.map((testcase, id) => (
    //                                 <tr key={id}>
    //                                     <td>
    //                                         <pre>{testcase.input}</pre>
    //                                     </td>
    //                                     <td>
    //                                         <pre>{testcase.execResult.output}</pre>
    //                                     </td>
    //                                     <td>
    //                                         <pre>{testcase.output}</pre>
    //                                     </td>
    //                                     <td className="text-center">
    //                                         {getVerdict(testcase.execResult)}
    //                                     </td>
    //                                     <td className="text-center">
    //                                         <pre>{testcase.execResult.runtime + ' ms'}</pre>
    //                                     </td>
    //                                     <td className="text-center">
    //                                         <pre>{testcase.execResult.memory + ' KB'}</pre>
    //                                     </td>
    //                                 </tr>
    //                             ))
    //                         }
    //                         </tbody>
    //                     </Table>
    //                 </div>
    //             </div>
    //         );
    //     }
    // };

    // const getTestcaseResults = () => {
    //     if (currentProblem && currentProblem.testcases && currentProblem.testcases.length > 0) {
    //         return (
    //             <Accordion style={{maxHeight: '72vh', overflowY: "auto", overflowX: "auto"}}>
    //                 {
    //                     currentProblem.testcases.map((testcase, id) => (
    //                         <Accordion.Item key={id} eventKey={id}>
    //                             <Accordion.Header>{`Test: #${id + 1}, time: ${testcase.execResult.runtime} ms., memory: ${testcase.execResult.memory} KB, verdict: ${testcase.execResult.verdict}`}</Accordion.Header>
    //                             <Accordion.Body style={{maxHeight: '65.6vh', overflowY: "auto", overflowX: "auto"}}>
    //                                 <Table id={id} bordered responsive="sm" size="sm">
    //                                     <thead style={{position: '', top: 0, zIndex: 1, background: '#fff'}}>
    //                                     <tr className="text-center">
    //                                         <th>INPUT</th>
    //                                         <th>OUTPUT</th>
    //                                         <th>ANSWER</th>
    //                                     </tr>
    //                                     </thead>
    //                                     <tbody>
    //                                         <tr key={id}>
    //                                             <td>
    //                                                 <pre>{testcase.input}</pre>
    //                                             </td>
    //                                             <td>
    //                                                 <pre>{testcase.execResult.output}</pre>
    //                                             </td>
    //                                             <td>
    //                                                 <pre>{testcase.output}</pre>
    //                                             </td>
    //                                         </tr>
    //                                     </tbody>
    //                                 </Table>
    //                             </Accordion.Body>
    //                         </Accordion.Item>
    //                     ))
    //                 }
    //             </Accordion>
    //         );
    //     }
    // }

    const getAlert = () => {
        if (!appState.config.workspaceDirectory || !appState.config.activeLanguage.lang || !appState.config.activeLanguage.compilationCommand || !appState.config.activeLanguage.compilationArgs) {
            return (
                <Row>
                    <Col>
                        <br/>
                        <Alert variant="danger" className="text-center">
                            Workspace directory is not set. Please go to Settings and set Workspace directory.
                        </Alert>
                    </Col>
                </Row>
            );
        }
    }

    return (
        <div>
            <Card body bg="light">
                <Row>
                    <Col xs={9}>
                        <Form.Group className="mb-3">
                            <Form.Control type="text" size="sm"
                                          placeholder="Enter Contest/Problem URL [Codeforces, AtCoder]"
                                          value={appState.url}
                                          disabled={!appState.config.workspaceDirectory}
                                          onChange={(e) => setAppState({...appState, url: e.target.value})}/>
                        </Form.Group>
                    </Col>
                    <Col xs={3}>
                        <div className="d-grid gap-2">
                            <Button size="sm" variant="outline-success" onClick={() => getProblemList()}
                                    disabled={disableActionButtons()}><FontAwesomeIcon
                                icon={faTasks}/> Load Problems</Button>
                        </div>
                    </Col>
                </Row>
                <Row>
                    <Col xs={6}>
                        <Form.Group className="mb-3">
                            <Form.Select size="sm" aria-label="Default select example"
                                         onChange={(e) => getSelectedProblemInfo(e.currentTarget.value)}>
                                {
                                    problemList.map((problem, id) => (
                                        <option key={id}
                                                value={problem.platform + '/' + problem.contestId + '/' + problem.label}>{problem.label.toUpperCase() + ' - ' + problem.name}
                                        </option>
                                    ))
                                }
                            </Form.Select>
                        </Form.Group>
                    </Col>
                    <Col xs={2}>
                        <div className="d-grid gap-2">
                            <Button size="sm" variant="outline-success" onClick={() => runTest()}
                                    disabled={!selectedProblem}><FontAwesomeIcon
                                icon={faTerminal}/> Test Code</Button>
                        </div>
                    </Col>
                    <Col xs={2}>
                        <div className="d-grid gap-2">
                            <Button size="sm" variant="outline-success" onClick={() => openSource()}
                                    disabled={!selectedProblem}><FontAwesomeIcon
                                icon={faFileCode}/> Open Code</Button>
                        </div>
                    </Col>
                    <Col xs={2}>
                        <div className="d-grid gap-2">
                            <Button size="sm" variant="outline-success" onClick={() => setShowCodeModal(true)}
                                    disabled={!selectedProblem}><FontAwesomeIcon
                                icon={faCode}/> View Code</Button>
                        </div>
                    </Col>
                </Row>
                <hr/>
                {/*<Row>*/}
                {/*    <Col xs={7}>*/}
                {/*        {currentProblem && currentProblem.metadata && (*/}
                {/*            <Form.Text>*/}
                {/*                <strong>*/}
                {/*                    /!*{currentProblem.metadata.label.toUpperCase() + ' - ' + currentProblem.metadata.name}*!/*/}
                {/*                    {'Time Limit: ' + currentProblem.metadata.timeLimit + ' sec, Memory Limit: ' + currentProblem.metadata.memoryLimit + ' MB'}*/}
                {/*                </strong>*/}
                {/*            </Form.Text>*/}
                {/*        )}*/}
                {/*    </Col>*/}
                {/*    <Col xs={5} style={{textAlign: 'right'}}>*/}
                {/*        <Form.Text>*/}
                {/*            {currentProblem && currentProblem.testStatus && getTestStatusText(currentProblem.testStatus)}*/}
                {/*        </Form.Text>*/}
                {/*    </Col>*/}
                {/*</Row>*/}
                {/*{getExecutionTable()}*/}
                {/*{getTestcaseResults()}*/}
                {getAlert()}
                {showCodeModal && (
                    <ViewCodeModal codePath={currentCodePath} setShowCodeModal={setShowCodeModal}/>
                )}
            </Card>
        </div>
    );
}