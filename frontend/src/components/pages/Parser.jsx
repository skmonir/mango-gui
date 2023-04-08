import { Alert, Button, Card, Col, Row, Spinner, Table } from "react-bootstrap";
import Form from "react-bootstrap/Form";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCheckCircle,
  faClock,
  faDownload,
  faFileCirclePlus,
  faSyncAlt,
  faTerminal,
  faTimesCircle
} from "@fortawesome/free-solid-svg-icons";
import { useEffect, useState } from "react";
import SocketClient from "../../socket/SocketClient.js";
import DataService from "../../services/DataService.js";
import Loading from "../misc/Loading.jsx";
import Utils from "../../Utils.js";
import AddCustomProblemModal from "../modals/AddCustomProblemModal.jsx";
import ShowToast from "../Toast/ShowToast.jsx";

function Parser({ appState }) {
  const socketClient = new SocketClient();

  const [parseUrl, setParseUrl] = useState("");
  const [initAlert, setInitAlert] = useState(false);
  const [parsingInProgress, setParsingInProgress] = useState(false);
  const [scheduler, setScheduler] = useState({
    inProgress: false,
    parsing: false,
    scheduled: false,
    url: "",
    startTime: ""
  });
  const [showToast, setShowToast] = useState(false);
  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: ""
  });
  const [parsingSingleProblem, setParsingSingleProblem] = useState(false);
  const [parsedProblemList, setParsedProblemList] = useState([]);
  const [showAddCustomProblemModal, setShowAddCustomProblemModal] = useState(
    false
  );

  useEffect(() => {
    fetchHistory();
    let socketConnParse = socketClient.initSocketConnection(
      "parse_problems_event",
      updateParseStatusFromSocket
    );
    let socketConnSchedule = socketClient.initSocketConnection(
      "parse_schedule_event",
      updateParseScheduleStatusFromSocket
    );
    return () => {
      socketConnParse.close();
      socketConnSchedule.close();
    };
  }, []);

  const fetchHistory = () => {
    DataService.getHistory().then(appHistory => {
      setParseUrl(appHistory?.parseUrl);
    });
  };

  const parseTriggerred = () => {
    setParsingInProgress(true);
    setTimeout(() => {
      DataService.parse(window.btoa(parseUrl)).then(data => {
        setParsedProblemList(data);
        setParsingInProgress(false);
        setInitAlert(true);
      });
    }, 0);
  };

  const parseSingleProblem = (index, url) => {
    setParsingInProgress(true);
    setParsingSingleProblem(true);
    setParsedProblemList(
      parsedProblemList.map((prob, i) =>
        i === index
          ? {
              ...prob,
              status: "running"
            }
          : prob
      )
    );
    console.log(parsedProblemList);
    DataService.parse(window.btoa(url)).then(data => {
      setParsedProblemList(
        parsedProblemList.map((prob, i) => (i === index ? data[0] : prob))
      );
      setParsingInProgress(false);
      setParsingSingleProblem(false);
    });
  };

  const scheduleParse = () => {
    setScheduler({ ...scheduler, inProgress: true });
    DataService.scheduleParse(window.btoa(parseUrl))
      .then(data => {
        console.log(data);
        setScheduler({
          ...scheduler,
          inProgress: false,
          scheduled: true,
          url: data.url,
          startTime: data.startTime
        });
      })
      .catch(error => {
        console.log(error);
        showToastMessage("Error", error.response.data.message);
        setScheduler({ ...scheduler, inProgress: false });
      });
  };

  const updateParseStatusFromSocket = data => {
    console.log(data);
    if (data.length > 1) {
      setParsedProblemList(data);
      setInitAlert(false);
      setParsingInProgress(false);
    }
  };

  const updateParseScheduleStatusFromSocket = data => {
    if (data.message === "done") {
      setScheduler({
        inProgress: false,
        parsing: false,
        scheduled: false,
        url: "",
        startTime: ""
      });
    } else if (data.message === "running") {
      setScheduler({
        ...scheduler,
        parsing: true
      });
    }
  };

  const createCustomProblem = () => {
    setShowAddCustomProblemModal(true);
  };

  const closeAddCustomProblemModal = () => {
    setShowAddCustomProblemModal(false);
  };

  const insertCustomProblemIntoList = data => {
    setParsedProblemList(data);
  };

  const getProblemStatusIcon = status => {
    if (!status || status === "running") {
      return <Spinner animation="border" variant="primary" size="sm" />;
    } else if (status === "failed") {
      return <FontAwesomeIcon style={{ color: "red" }} icon={faTimesCircle} />;
    } else if (status === "success") {
      return (
        <FontAwesomeIcon style={{ color: "green" }} icon={faCheckCircle} />
      );
    }
  };

  const showToastMessage = (variant, message) => {
    setShowToast(true);
    setToastMsgObj({
      variant: variant,
      message: message
    });
  };

  const disableActionButtons = () => {
    return (
      parsingInProgress ||
      scheduler.inProgress ||
      !appState.config.workspaceDirectory
    );
  };

  const getParsingTable = () => {
    return (
      <Table bordered striped responsive="sm" size="sm">
        <thead>
          <tr className="text-center">
            <th>#</th>
            <th style={{ minWidth: "60vh", maxWidth: "60vh" }}>Problem Name</th>
            <th>Status</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          {parsedProblemList.map((problem, id) => (
            <tr key={id}>
              <td className="text-center">
                <a
                  href="#"
                  onClick={() =>
                    DataService.openResource({ path: problem?.url })
                  }
                >
                  {problem?.label}
                </a>
              </td>
              <td>
                <a
                  href="#"
                  onClick={() =>
                    DataService.openResource({ path: problem?.url })
                  }
                >
                  {problem?.name}
                </a>
              </td>
              <td className="text-center">
                {getProblemStatusIcon(problem?.status)}
              </td>
              <td className="text-center">
                <Button
                  variant="outline-success"
                  size="sm"
                  onClick={() => parseSingleProblem(id, problem?.url)}
                  disabled={
                    parsingInProgress ||
                    !appState.config.workspaceDirectory ||
                    problem?.url.startsWith("custom/")
                  }
                >
                  <FontAwesomeIcon icon={faSyncAlt} /> Refresh
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    );
  };

  const getParserBody = () => {
    if (parsingInProgress && !parsingSingleProblem) {
      return <Loading />;
    } else if (parsedProblemList && parsedProblemList.length > 0) {
      return getParsingTable();
    } else if (initAlert) {
      return (
        <Alert variant="warning" className="text-center p-1 mb-2">
          Oops! Something went wrong! Please <strong>check the URL</strong> or
          try again.
        </Alert>
      );
    }
  };

  return (
    <div>
      <div className="panel">
        <div className="panel-body">
          <Row>
            <Col xs={5}>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="Enter Contest/Problem URL [Codeforces, AtCoder]"
                value={parseUrl}
                disabled={!appState.config.workspaceDirectory}
                onChange={e => setParseUrl(e.target.value)}
              />
            </Col>
            <Col xs={2}>
              <div className="d-grid gap-2">
                <Button
                  size="sm"
                  variant="outline-primary"
                  onClick={() => parseTriggerred()}
                  disabled={
                    disableActionButtons() || Utils.isStrNullOrEmpty(parseUrl)
                  }
                >
                  {parsingInProgress ? (
                    <Spinner
                      as="span"
                      animation="border"
                      size="sm"
                      role="status"
                      aria-hidden="true"
                    />
                  ) : (
                    <FontAwesomeIcon icon={faDownload} />
                  )}{" "}
                  {" Parse Testcases"}
                </Button>
              </div>
            </Col>
            <Col xs={2}>
              <div className="d-grid gap-2">
                <Button
                  size="sm"
                  variant="outline-success"
                  onClick={scheduleParse}
                  disabled={
                    disableActionButtons() ||
                    Utils.isStrNullOrEmpty(parseUrl) ||
                    scheduler?.scheduled
                  }
                >
                  {scheduler?.inProgress ? (
                    <Spinner
                      as="span"
                      animation="border"
                      size="sm"
                      role="status"
                      aria-hidden="true"
                    />
                  ) : (
                    <FontAwesomeIcon icon={faClock} />
                  )}{" "}
                  {` Schedule${scheduler.scheduled ? "d" : ""}`}
                </Button>
              </div>
            </Col>
            <Col xs={3}>
              <div className="d-grid gap-2">
                <Button
                  size="sm"
                  variant="outline-success"
                  onClick={() => createCustomProblem()}
                  disabled={disableActionButtons()}
                >
                  <FontAwesomeIcon icon={faFileCirclePlus} /> Create Custom
                  Problem
                </Button>
              </div>
            </Col>
          </Row>
          <hr />
          {scheduler?.scheduled && (
            <Row>
              <Col xs={12}>
                <Alert variant="success" className="text-center p-1 mb-2">
                  <span>
                    {scheduler.parsing && (
                      <Spinner
                        as="span"
                        animation="border"
                        size="sm"
                        role="status"
                        aria-hidden="true"
                      />
                    )}
                  </span>{" "}
                  <strong>{scheduler?.url}</strong> is scheduled to be parsed on{" "}
                  <strong>
                    {Utils.dateStringToUiFormat(scheduler?.startTime)}
                  </strong>
                </Alert>
              </Col>
            </Row>
          )}
          {getParserBody()}
          {!appState.config.workspaceDirectory && (
            <Row>
              <Col>
                <br />
                <Alert variant="danger" className="text-center p-1 mb-2">
                  Configuration is not set properly. Please go to Settings and
                  set necessary fields.
                </Alert>
              </Col>
            </Row>
          )}
        </div>
      </div>
      {showAddCustomProblemModal && (
        <AddCustomProblemModal
          closeAddCustomProblemModal={closeAddCustomProblemModal}
          insertCustomProblemIntoList={insertCustomProblemIntoList}
        />
      )}
      {showToast && (
        <ShowToast toastMsgObj={toastMsgObj} setShowToast={setShowToast} />
      )}
    </div>
  );
}

export default Parser;
