import {Alert, Button, Card, Col, Row, Spinner, Table} from "react-bootstrap";
import Form from "react-bootstrap/Form";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCheckCircle, faDownload, faSyncAlt, faTimesCircle} from "@fortawesome/free-solid-svg-icons";
import {useEffect, useState} from "react";
import SocketClient from "../socket/SocketClient.js";
import DataService from "../services/DataService.js";
import Loading from "./Loading.jsx";

function Parser({appState, setAppState}) {
    const socketClient = new SocketClient();

    const [initAlert, setInitAlert] = useState(false);
    const [parsingInProgress, setParsingInProgress] = useState(false);
    const [parsingSingleProblem, setParsingSingleProblem] = useState(false);
    const [parsedProblemList, setParsedProblemList] = useState([]);

    useEffect(() => {
        let socketConn = socketClient.initSocketConnection("parse_problems_event", updateParseStatusFromSocket);
        return () => {
            socketConn.close()
        }
    }, []);

    const parseTriggerred = () => {
        setParsingInProgress(true);
        setTimeout(() => {
            DataService.parse(window.btoa(appState.url)).then(data => {
                setParsedProblemList(data);
                setAppState({...appState, parsedProblemList: data});
                setParsingInProgress(false);
                setInitAlert(true);
            });
        }, 0);
    }

    const parseSingleProblem = (index, url) => {
        setParsingInProgress(true);
        setParsingSingleProblem(true);
        setParsedProblemList(parsedProblemList.map((prob, i) => i === index ? {
            ...prob,
            parseStatus: 'running'
        } : prob));
        DataService.parse(window.btoa(url)).then(data => {
            setParsedProblemList(parsedProblemList.map((prob, i) => i === index ? data[0] : prob));
            setAppState({...appState, parsedProblemList: [...parsedProblemList]});
            setParsingInProgress(false);
            setParsingSingleProblem(false);
        });
    };

    const updateParseStatusFromSocket = (data) => {
        console.log(data);
        if (data.length > 1) {
            setParsedProblemList(data);
            setParsingInProgress(false);
        }
    }

    const getProblemStatusIcon = (status) => {
        if (!status || status === 'running') {
            return <Spinner animation="border" variant="primary" size="sm"/>;
        } else if (status === 'failed') {
            return <FontAwesomeIcon style={{color: 'red'}} icon={faTimesCircle}/>;
        } else if (status === 'success') {
            return <FontAwesomeIcon style={{color: 'green'}} icon={faCheckCircle}/>
        }
    };

    const disableActionButtons = () => {
        return !appState.url || appState.url === '' || parsingInProgress || !appState.config.workspaceDirectory;
    }

    const getParsingTable = () => {
        return (
            <Table bordered striped responsive="sm" size="sm">
                <thead>
                <tr className="text-center">
                    <th>#</th>
                    <th style={{minWidth: '60vh', maxWidth: '60vh'}}>Problem Name</th>
                    <th>Status</th>
                    <th>Action</th>
                </tr>
                </thead>
                <tbody>
                {
                    parsedProblemList.map((problem, id) => (
                        <tr key={id}>
                            <td className="text-center"><a href={problem?.url}
                                                           style={{pointerEvents: "none"}}>{problem?.label.toUpperCase()}</a>
                            </td>
                            <td><a href={problem?.url}
                                   style={{pointerEvents: "none"}}>{problem?.name}</a></td>
                            <td className="text-center">{getProblemStatusIcon(problem?.status)}</td>
                            <td className="text-center">
                                <Button variant="outline-success" size="sm"
                                        onClick={() => parseSingleProblem(id, problem?.url)}
                                        disabled={parsingInProgress || !appState.config.workspaceDirectory}>
                                    <FontAwesomeIcon icon={faSyncAlt}/> Refresh
                                </Button>
                            </td>
                        </tr>
                    ))
                }
                </tbody>
            </Table>
        );
    };

    const getParserBody = () => {
        if (parsingInProgress && !parsingSingleProblem) {
            return <Loading/>;
        } else if (parsedProblemList && parsedProblemList.length > 0) {
            return getParsingTable();
        } else if (initAlert) {
            return (
                <Alert variant="warning" className="text-center">
                    Oops! Something went wrong! Please <strong>check the URL</strong> or try again.
                </Alert>
            );
        }
    }

    return (
        <div>
            <Card body bg="light">
                <Row>
                    <Col xs={9}>
                        <Form.Control type="text" size="sm" placeholder="Enter Contest/Problem URL [Codeforces, AtCoder]"
                                      value={appState.url}
                                      disabled={!appState.config.workspaceDirectory}
                                      onChange={(e) => setAppState({...appState, url: e.target.value})}/>
                    </Col>
                    <Col>
                        <div className="d-grid gap-2">
                            <Button size="sm" variant="outline-success" onClick={() => parseTriggerred()}
                                    disabled={disableActionButtons()}><FontAwesomeIcon
                                icon={faDownload}/> Parse
                                Testcases</Button>
                        </div>
                    </Col>
                </Row>
                <hr/>
                {
                    getParserBody()
                }
                {!appState.config.workspaceDirectory && (
                    <Row>
                        <Col>
                            <br/>
                            <Alert variant="danger" className="text-center">
                                Workspace directory is not set. Please go to Settings and set Workspace directory.
                            </Alert>
                        </Col>
                    </Row>
                )}
            </Card>
        </div>
    );
}

export default Parser