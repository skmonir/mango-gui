import { Card, Tab, Tabs } from "react-bootstrap";
import InputGenerator from "./InputGenerator.jsx";

export default function TestcaseGenerator({ appState }) {
  return (
    <Card body bg="light">
      <Tabs defaultActiveKey="input_generator" className="mb-3">
        <Tab eventKey="input_generator" title="Input Generator">
          <InputGenerator appState={appState} />
        </Tab>
        <Tab eventKey="output_generator" title="Output Generator"></Tab>
      </Tabs>
    </Card>
  );
}
