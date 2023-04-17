import { Alert, Button, Card, Col, Row, Spinner, Table } from "react-bootstrap";
import Form from "react-bootstrap/Form";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBan,
  faCheckCircle,
  faClock,
  faDownload,
  faFileCirclePlus,
  faRefresh,
  faSyncAlt,
  faTimesCircle,
} from "@fortawesome/free-solid-svg-icons";
import { useEffect, useState } from "react";
import SocketClient from "../../socket/SocketClient.js";
import DataService from "../../services/DataService.js";
import Loading from "../misc/Loading.jsx";
import Utils from "../../Utils.js";
import AddCustomProblemModal from "../modals/AddCustomProblemModal.jsx";
import ShowToast from "../Toast/ShowToast.jsx";
import { confirmDialog } from "../modals/ConfirmationDialog.jsx";

function Parser({ config, appData }) {
  const socketClient = new SocketClient();

  const [parseUrl, setParseUrl] = useState("");
  const [initAlert, setInitAlert] = useState(false);
  const [parseSchedulerTasks, setParseSchedulerTasks] = useState([]);
  const [showToast, setShowToast] = useState(false);
  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: "",
  });
  const [parsedProblemList, setParsedProblemList] = useState([]);
  const [showAddCustomProblemModal, setShowAddCustomProblemModal] =
    useState(false);
  const [ipFlags, setIpFlags] = useState({
    parsingWithUrl: false,
    parsingWithRefresh: false,
    scheduling: false,
    refreshingScheduledTasks: false,
  });

  useEffect(() => {
    setParseUrl(appData?.queryHistories?.parseUrl);
    setParseSchedulerTasks(appData?.parseSchedulerTasks);

    let socketConnParse = socketClient.initSocketConnection(
      "parse_problems_event",
      updateParseStatusFromSocket
    );
    let socketConnSchedule = socketClient.initSocketConnection(
      "parse_schedule_event",
      updateParseScheduledTasksFromSocket
    );
    return () => {
      socketConnParse.close();
      socketConnSchedule.close();
    };
  }, []);

  const parseTriggerred = () => {
    setIpFlags({
      ...ipFlags,
      parsingWithUrl: true,
    });
    setTimeout(() => {
      DataService.parse(window.btoa(parseUrl)).then((data) => {
        setParsedProblemList(data);
        setIpFlags({
          ...ipFlags,
          parsingWithUrl: false,
        });
        setInitAlert(true);
      });
    }, 0);
  };

  const parseSingleProblem = (index, url) => {
    setIpFlags({
      ...ipFlags,
      parsingWithUrl: true,
      parsingWithRefresh: true,
    });
    setParsedProblemList(
      parsedProblemList.map((prob, i) =>
        i === index
          ? {
              ...prob,
              status: "running",
            }
          : prob
      )
    );
    console.log(parsedProblemList);
    DataService.parse(window.btoa(url)).then((data) => {
      setParsedProblemList(
        parsedProblemList.map((prob, i) => (i === index ? data[0] : prob))
      );
      setIpFlags({
        ...ipFlags,
        parsingWithUrl: false,
        parsingWithRefresh: false,
      });
    });
  };

  const scheduleParse = () => {
    setIpFlags({
      ...ipFlags,
      scheduling: true,
    });
    DataService.scheduleParse({ url: parseUrl })
      .then((data) => {
        console.log(data);
      })
      .catch((error) => {
        console.log(error);
        showToastMessage("Error", error.response.data.message);
      })
      .finally(() =>
        setIpFlags({
          ...ipFlags,
          scheduling: false,
        })
      );
  };

  const removeScheduledTaskTriggered = (taskId) => {
    confirmDialog({
      title: "Delete Confirmation!",
      message: "Are you sure to delete this scheduled task?",
      okButton: {
        label: "Yes, Delete!",
        variant: "outline-danger",
      },
    }).then((response) => {
      if (response?.ok) {
        removeScheduledTask(taskId);
      }
    });
  };

  const removeScheduledTask = (taskId) => {
    DataService.removeParseScheduledTask(taskId).then((data) => {
      console.log(data);
    });
  };

  const updateParseStatusFromSocket = (data) => {
    console.log(data);
    if (data.length > 1) {
      setParsedProblemList(data);
      setInitAlert(false);
      setIpFlags({
        ...ipFlags,
        parsingWithUrl: false,
      });
    }
  };

  const updateParseScheduledTasksFromSocket = (tasks) => {
    console.log(tasks);
    setParseSchedulerTasks(tasks);
  };

  const getScheduledTasks = () => {
    setIpFlags({
      ...ipFlags,
      refreshingScheduledTasks: true,
    });
    setTimeout(() => {
      DataService.getParseScheduledTasks()
        .then((tasks) => {
          setParseSchedulerTasks(tasks);
        })
        .finally(() => {
          setIpFlags({
            ...ipFlags,
            refreshingScheduledTasks: false,
          });
        });
    }, 700);
  };

  const createCustomProblem = () => {
    setShowAddCustomProblemModal(true);
  };

  const closeAddCustomProblemModal = () => {
    setShowAddCustomProblemModal(false);
  };

  const insertCustomProblemIntoList = (data) => {
    setParsedProblemList(data);
  };

  const getProblemStatusIcon = (status) => {
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
      message: message,
    });
  };

  const disableActionButtons = () => {
    let disable = false;
    for (const [_, value] of Object.entries(ipFlags)) {
      disable = disable || value;
    }
    return disable || !config.workspaceDirectory;
  };

  const getSchedulerRowColor = (stage) => {
    if (
      stage === "SCHEDULED" ||
      stage === "RE_SCHEDULED" ||
      stage === "COMPLETE"
    ) {
      return "table-success";
    } else if (stage === "RUNNING") {
      return "table-warning";
    } else {
      return "table-danger";
    }
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
                <pre className="mb-0">
                  <a
                    href="#"
                    onClick={() =>
                      DataService.openResource({ path: problem?.url })
                    }
                  >
                    {problem?.label}
                  </a>
                </pre>
              </td>
              <td>
                <pre className="mb-0">
                  <a
                    href="#"
                    onClick={() =>
                      DataService.openResource({ path: problem?.url })
                    }
                  >
                    {problem?.name}
                  </a>
                </pre>
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
                    ipFlags.parsingWithUrl ||
                    !config.workspaceDirectory ||
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

  const getParseSchedulerTasksTable = () => {
    return (
      <Table bordered responsive="sm" size="sm">
        <thead>
          <tr className="text-center">
            <th style={{ minWidth: "60vh", maxWidth: "60vh" }}>Parsing URL</th>
            <th>Scheduled Time</th>
            <th>Stage</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          {parseSchedulerTasks.map((scheduledTask) => (
            <tr
              key={scheduledTask.id}
              className={getSchedulerRowColor(scheduledTask.stage)}
            >
              <td>
                <pre className="mb-0">
                  <a
                    href="#"
                    onClick={() =>
                      DataService.openResource({ path: scheduledTask.url })
                    }
                  >
                    {scheduledTask.url}
                  </a>
                </pre>{" "}
                {scheduledTask?.stage === "RUNNING" && (
                  <span>
                    <Spinner
                      as="span"
                      animation="border"
                      size="sm"
                      role="status"
                      aria-hidden="true"
                      variant="success"
                    />
                  </span>
                )}
              </td>
              <td className="text-center">
                <pre className="mb-0">
                  {Utils.dateToLocaleString(scheduledTask.startTime)}
                </pre>
              </td>
              <td className="text-center">
                <pre className="mb-0">{scheduledTask?.stage}</pre>
              </td>
              <td className="text-center">
                {(scheduledTask.stage === "SCHEDULED" ||
                  scheduledTask.stage === "RE_SCHEDULED") && (
                  <Button
                    variant="outline-danger"
                    size="sm"
                    disabled={disableActionButtons()}
                    onClick={() =>
                      removeScheduledTaskTriggered(scheduledTask.id)
                    }
                  >
                    <FontAwesomeIcon icon={faBan} /> Delete Task
                  </Button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    );
  };

  const getParserBody = () => {
    if (ipFlags.parsingWithUrl && !ipFlags.parsingWithRefresh) {
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
                disabled={!config.workspaceDirectory}
                onChange={(e) => setParseUrl(e.target.value)}
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
                  {ipFlags.parsingWithUrl ? (
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
                    disableActionButtons() || Utils.isStrNullOrEmpty(parseUrl)
                  }
                >
                  {ipFlags.scheduling ? (
                    <Spinner
                      as="span"
                      animation="border"
                      size="sm"
                      role="status"
                      aria-hidden="true"
                    />
                  ) : (
                    <FontAwesomeIcon icon={faClock} />
                  )}
                  {" Schedule"}
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
          {parseSchedulerTasks?.length > 0 && (
            <>
              <Row className="mt-0">
                <Col xs="10"></Col>
                <Col xs="2" className="d-grid gap-2">
                  <Button
                    size="sm"
                    variant="success"
                    className="mb-1"
                    disabled={disableActionButtons()}
                    onClick={getScheduledTasks}
                  >
                    {ipFlags.refreshingScheduledTasks ? (
                      <Spinner
                        as="span"
                        animation="border"
                        size="sm"
                        role="status"
                        aria-hidden="true"
                        variant="warning"
                      />
                    ) : (
                      <FontAwesomeIcon icon={faRefresh} />
                    )}{" "}
                    Refresh
                  </Button>
                </Col>
              </Row>
              <Row>
                <Col xs={12}>{getParseSchedulerTasksTable()}</Col>
              </Row>
            </>
          )}
          {getParserBody()}
          {!config.workspaceDirectory && (
            <Row>
              <Col>
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
