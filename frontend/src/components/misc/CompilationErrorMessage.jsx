import { Table } from "react-bootstrap";

export default function CompilationErrorMessage({ props }) {
  return (
    <div
      style={{
        minHeight: props.minHeight,
        maxHeight: props.maxHeight,
        overflowY: "auto",
        overflowX: "auto",
        border: "2px solid transparent",
        borderColor: "black",
        borderRadius: "5px",
      }}
    >
      <Table bordered responsive="sm" size="sm">
        <tbody>
          <tr>
            <td className="table-danger">
              <pre>{props.error}</pre>
            </td>
          </tr>
        </tbody>
      </Table>
    </div>
  );
}
