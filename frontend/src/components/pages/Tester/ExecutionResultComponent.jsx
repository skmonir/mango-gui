import { Button, ButtonGroup, Col, Row, Spinner } from "react-bootstrap";
import CompilationErrorMessage from "../../misc/CompilationErrorMessage.jsx";
import ReactTable from "react-table";
import AC from "../../../assets/icons/AC.svg";
import WA from "../../../assets/icons/WA.svg";
import CE from "../../../assets/icons/CE.svg";
import RE from "../../../assets/icons/RE.svg";
import TLE from "../../../assets/icons/TLE.svg";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faClone, faEdit, faTrashAlt } from "@fortawesome/free-solid-svg-icons";

export default function ExecutionResultComponent({
  selectedProblemFilteredExecResult,
  disableActionButtons,
  cloneUpdateCustomTest,
  deleteCustomTestTriggered,
}) {
  const verdictIcons = {
    AC: AC,
    WA: WA,
    CE: CE,
    RE: RE,
    TLE: TLE,
  };

  const getVerdict = (testcaseExecutionDetails) => {
    if (testcaseExecutionDetails?.status === "running") {
      return (
        <div className="d-flex justify-content-center">
          <Spinner animation="border" variant="primary" size="sm" />
        </div>
      );
    } else if (testcaseExecutionDetails?.status !== "none") {
      return (
        <pre
          style={{
            textAlign: "center",
            color:
              testcaseExecutionDetails?.testcaseExecutionResult?.verdict ===
              "AC"
                ? "green"
                : "red",
          }}
        >
          <img
            src={
              verdictIcons[
                testcaseExecutionDetails?.testcaseExecutionResult?.verdict
              ]
            }
            style={{ maxWidth: "30px" }}
          />{" "}
          <strong>
            {testcaseExecutionDetails?.testcaseExecutionResult?.verdict}
          </strong>
        </pre>
      );
    }
  };

  const getTestcaseRowColor = (testcaseExecutionDetails) => {
    if (["none", "running"].includes(testcaseExecutionDetails?.status)) {
      return "#e2e3e5";
    } else {
      if (testcaseExecutionDetails?.testcaseExecutionResult?.verdict === "AC") {
        return "#d1e7dd";
      } else {
        return "#f8d7da";
      }
    }
  };

  const getExecTableActionButtons = (execDetails) => {
    return (
      <ButtonGroup className="d-flex justify-content-center">
        <Button
          size="sm"
          variant="primary"
          title="Edit"
          onClick={() => cloneUpdateCustomTest(execDetails.testcase, "Update")}
          disabled={disableActionButtons()}
        >
          <FontAwesomeIcon icon={faEdit} />
        </Button>
        <Button
          size="sm"
          variant="success"
          title="Clone"
          onClick={() => cloneUpdateCustomTest(execDetails.testcase, "Clone")}
          disabled={disableActionButtons()}
        >
          <FontAwesomeIcon icon={faClone} />
        </Button>
        <Button
          size="sm"
          variant="danger"
          title="Delete"
          onClick={() =>
            deleteCustomTestTriggered(execDetails.testcase.inputFilePath)
          }
          disabled={disableActionButtons()}
        >
          <FontAwesomeIcon icon={faTrashAlt} />
        </Button>
      </ButtonGroup>
    );
  };

  const getCompactExecutionTable = () => {
    const execTableColumns = [
      {
        Header: () => (
          <span style={{ textAlign: "center" }}>
            <strong>Testcase</strong>
          </span>
        ),
        accessor: "inputFilePath",
        Cell: ({ original }) => {
          return (
            <pre
              style={{
                textAlign: "center",
                overflow: "hidden",
                textOverflow: "ellipsis",
              }}
            >
              {original?.testcase?.inputFilePath
                .split("\\")
                .pop()
                .split("/")
                .pop()}
            </pre>
          );
        },
      },
      {
        Header: () => (
          <span style={{ textAlign: "center" }}>
            <strong>Result</strong>
          </span>
        ),
        accessor: "event",
        maxWidth: 75,
        Cell: ({ original }) => {
          return getVerdict(original);
        },
      },
      {
        Header: () => (
          <span style={{ textAlign: "center" }}>
            <strong>Time</strong>
          </span>
        ),
        accessor: "event",
        maxWidth: 80,
        Cell: ({ original }) => {
          return (
            <pre style={{ textAlign: "center" }}>
              {original?.testcaseExecutionResult?.consumedTime + " ms"}
            </pre>
          );
        },
      },
      {
        Header: () => (
          <span style={{ textAlign: "center" }}>
            <strong>Memory</strong>
          </span>
        ),
        accessor: "event",
        maxWidth: 85,
        Cell: ({ original }) => {
          return (
            <pre style={{ textAlign: "center" }}>
              {original?.testcaseExecutionResult?.consumedMemory + " KB"}
            </pre>
          );
        },
      },
      {
        Header: () => (
          <span style={{ textAlign: "center" }}>
            <strong>Actions</strong>
          </span>
        ),
        accessor: "event",
        maxWidth: 120,
        Cell: ({ original }) => {
          return getExecTableActionButtons(original);
        },
      },
    ];

    const ioTableColumns = [
      {
        Header: () => (
          <span style={{ textAlign: "center" }}>
            <strong>Input</strong>
          </span>
        ),
        accessor: "input",
        minWidth: 300,
        Cell: ({ original }) => {
          return (
            <pre style={{ overflow: "hidden", textOverflow: "ellipsis" }}>
              {original?.testcase?.input}
            </pre>
          );
        },
      },
      {
        Header: () => (
          <span style={{ textAlign: "center" }}>
            <strong>Expected Output</strong>
          </span>
        ),
        accessor: "expectedOutput",
        minWidth: 375,
        Cell: ({ original }) => {
          return (
            <pre style={{ overflow: "hidden", textOverflow: "ellipsis" }}>
              {original?.testcase?.output}
            </pre>
          );
        },
      },
      {
        Header: () => (
          <span style={{ textAlign: "center" }}>
            <strong>Program Output</strong>
          </span>
        ),
        accessor: "programOutput",
        minWidth: 375,
        Cell: ({ original }) => {
          return (
            <pre style={{ overflow: "hidden", textOverflow: "ellipsis" }}>
              {original?.testcaseExecutionResult?.output}
            </pre>
          );
        },
      },
    ];

    return (
      <div>
        <ReactTable
          data={selectedProblemFilteredExecResult.testcaseExecutionDetailsList}
          data-testid="exec_result_table"
          columns={execTableColumns}
          sortable={false}
          showPagination={false}
          showPageSizeOptions={false}
          resizable={false}
          collapseOnDataChange={false}
          minRows={0}
          defaultPageSize={1000}
          SubComponent={(rowInfo) => {
            return (
              <div>
                <ReactTable
                  data={[rowInfo.original]}
                  data-testid={`io_table_${rowInfo.index + 1}`}
                  columns={ioTableColumns}
                  sortable={false}
                  showPagination={false}
                  showPageSizeOptions={false}
                  minRows={0}
                  className="-bordered"
                />
              </div>
            );
          }}
          getTrProps={(state, rowInfo, instance) => {
            if (rowInfo === undefined) {
              return {};
            }
            return {
              style: {
                background: getTestcaseRowColor(rowInfo.original),
              },
            };
          }}
        />
      </div>
    );
  };

  return (
    <>
      {selectedProblemFilteredExecResult?.compilationError && (
        <Row className="mb-2">
          <Col xs={12}>
            <CompilationErrorMessage
              props={{
                maxHeight: "35vh",
                error: selectedProblemFilteredExecResult?.compilationError,
              }}
            />
          </Col>
        </Row>
      )}
      <Row>
        {selectedProblemFilteredExecResult?.testcaseExecutionDetailsList &&
          getCompactExecutionTable()}
      </Row>
    </>
  );
}
