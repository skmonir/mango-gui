import { Table } from "react-bootstrap";

export default function GeneratorLogs({ props }) {
  return (
    <div
      style={{
        minHeight: props.minHeight,
        maxHeight: props.maxHeight,
        overflowY: "auto",
        overflowX: "auto",
        border: "2px solid transparent",
        borderColor: "black",
        borderRadius: "5px"
      }}
    >
      <Table bordered responsive="sm" size="sm">
        <tbody>
          {props.logList
            .filter(log => log.status === "success")
            .slice(0)
            .reverse()
            .map((t, id) => (
              <tr
                key={id}
                className={
                  t.testcaseExecutionResult.executionError !== ""
                    ? "table-danger"
                    : "table-success"
                }
              >
                <td>
                  <pre>{t.testcase.execOutputFilePath}</pre>
                </td>
              </tr>
            ))}
        </tbody>
      </Table>
    </div>
  );
}
