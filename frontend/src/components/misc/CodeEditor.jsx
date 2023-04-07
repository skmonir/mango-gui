import { Button, Col, Modal, Row } from "react-bootstrap";
import Form from "react-bootstrap/Form";
import { useEffect, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCog, faSave } from "@fortawesome/free-solid-svg-icons";

import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-java";
import "ace-builds/src-noconflict/mode-python";

import "ace-builds/src-noconflict/snippets/java";
import "ace-builds/src-noconflict/snippets/python";

import "ace-builds/src-noconflict/theme-monokai";
import "ace-builds/src-noconflict/theme-xcode";
import "ace-builds/src-noconflict/theme-textmate";
import "ace-builds/src-noconflict/theme-twilight";
import "ace-builds/src-noconflict/theme-terminal";

import "ace-builds/src-noconflict/ext-language_tools";
import "ace-builds/src-min-noconflict/ext-searchbox";
import DataService from "../../services/DataService.js";
import Utils from "../../Utils.js";

export default function CodeEditor({
  codeContent,
  onChange,
  onBlur,
  readOnly,
  customElementsOnHeader
}) {
  const themes = ["monokai", "xcode", "textmate", "twilight", "terminal"];
  const fontSizes = ["13", "14", "16", "18", "20", "24", "28", "32", "40"];
  const [editorPref, setEditorPref] = useState({
    theme: "monokai",
    fontSize: "14",
    tabSize: "4"
  });
  const [modalSetup, setModalSetup] = useState({});
  const [editorLang, setEditorLang] = useState("");

  const [showEditorPref, setShowEditorPref] = useState(false);

  useEffect(() => {
    if (codeContent.lang !== "tgen") {
      getEditorPreference();
    }
    convertModeToLang();
  }, []);

  const getEditorPreference = () => {
    DataService.getEditorPreference().then(pref => {
      handleEditorPrefResponse(pref);
    });
  };

  const updateEditorPreference = preference => {
    DataService.updateEditorPreference(preference).then(pref => {
      handleEditorPrefResponse(pref);
    });
  };

  const handleEditorPrefResponse = pref => {
    setEditorPref({
      ...editorPref,
      theme: Utils.isStrNullOrEmpty(pref.theme) ? "monokai" : pref.theme,
      fontSize: Utils.isStrNullOrEmpty(pref.fontSize) ? "14" : pref.fontSize,
      tabSize: Utils.isStrNullOrEmpty(pref.tabSize) ? "4" : pref.tabSize
    });
  };

  const convertModeToLang = () => {
    if (
      codeContent.lang === "" ||
      codeContent.lang === "tgen" ||
      codeContent.lang === "cpp"
    ) {
      setEditorLang("java");
    } else {
      setEditorLang(codeContent.lang);
    }
  };

  const openEditorPrefModal = () => {
    console.log(editorPref);
    setModalSetup({
      ...editorPref
    });
    setTimeout(() => setShowEditorPref(true), 300);
  };

  const saveEditorPref = () => {
    console.log(modalSetup);
    updateEditorPreference(modalSetup);
    setShowEditorPref(false);
  };

  return (
    <div>
      <Row style={{ marginTop: "3px", marginBottom: "5px" }}>
        {!readOnly.preference && (
          <Col xs="auto">
            <FontAwesomeIcon
              icon={faCog}
              onClick={openEditorPrefModal}
              style={{ cursor: "pointer", marginTop: "8px" }}
            />
          </Col>
        )}
        <Col xs="auto">
          <Form.Select size="sm" value={codeContent.lang} disabled={true}>
            <option value=""></option>
            <option value="cpp">CPP</option>
            <option value="java">Java</option>
            <option value="python">Python</option>
            <option value="tgen">TGen</option>
          </Form.Select>
        </Col>
        {customElementsOnHeader &&
          customElementsOnHeader.map((elem, idx) => (
            <Col key={idx} xs={elem.colSpan}>
              {elem.elem}
            </Col>
          ))}
      </Row>
      <Row>
        <Col xs={12} style={{ minHeight: "86vh", overflowY: "auto" }}>
          <AceEditor
            height={"100%"}
            width={"100%"}
            mode={editorLang}
            theme={editorPref.theme}
            name="code_editor"
            onChange={code => onChange(code)}
            onBlur={onBlur}
            showPrintMargin={false}
            showGutter={true}
            highlightActiveLine={true}
            value={codeContent.code}
            setOptions={{
              readOnly: readOnly.editor,
              enableBasicAutocompletion: true,
              enableLiveAutocompletion: true,
              enableSnippets: true,
              showLineNumbers: true,
              dragEnabled: true,
              tabSize: Number(editorPref.tabSize),
              fontSize: Number(editorPref.fontSize)
            }}
          />
        </Col>
      </Row>
      <Modal
        show={showEditorPref}
        onHide={() => setShowEditorPref(false)}
        size="sm"
        aria-labelledby="contained-modal-title-vcenter"
        centered
      >
        <Modal.Header
          style={{ paddingBottom: "5px", paddingTop: "5px" }}
          closeButton
        >
          <strong>Editor Preferences</strong>
        </Modal.Header>
        <Modal.Body style={{ paddingBottom: "2px", paddingTop: "5px" }}>
          <Row>
            <Col xs={12}>
              <Form.Group className="mb-3">
                <Form.Label>Theme</Form.Label>
                <Form.Select
                  size="sm"
                  aria-label="Default select example"
                  value={modalSetup.theme}
                  onChange={e =>
                    setModalSetup({
                      ...modalSetup,
                      theme: e.currentTarget.value
                    })
                  }
                >
                  {themes.map((theme, id) => (
                    <option key={id} value={theme}>
                      {theme}
                    </option>
                  ))}
                </Form.Select>
              </Form.Group>
            </Col>
          </Row>
          <Row>
            <Col xs={12}>
              <Form.Group className="mb-3">
                <Form.Label>Font Size</Form.Label>
                <Form.Select
                  size="sm"
                  aria-label="Default select example"
                  value={modalSetup.fontSize}
                  onChange={e =>
                    setModalSetup({
                      ...modalSetup,
                      fontSize: e.currentTarget.value
                    })
                  }
                >
                  {fontSizes.map((size, id) => (
                    <option key={id} value={size}>
                      {size}
                    </option>
                  ))}
                </Form.Select>
              </Form.Group>
            </Col>
          </Row>
          <Row>
            <Col xs={12}>
              <Form.Group className="mb-3">
                <Form.Label>Tab Size</Form.Label>
                <Form.Select
                  size="sm"
                  aria-label="Default select example"
                  value={modalSetup.tabSize}
                  onChange={e =>
                    setModalSetup({
                      ...modalSetup,
                      tabSize: e.currentTarget.value
                    })
                  }
                >
                  <option value="2">2</option>
                  <option value="4">4</option>
                  <option value="8">8</option>
                </Form.Select>
              </Form.Group>
            </Col>
          </Row>
        </Modal.Body>
        <Modal.Footer style={{ paddingBottom: "5px", paddingTop: "5px" }}>
          <Button size="sm" variant="outline-success" onClick={saveEditorPref}>
            <FontAwesomeIcon icon={faSave} /> Save
          </Button>
          <Button
            size="sm"
            variant="outline-secondary"
            onClick={() => setShowEditorPref(false)}
          >
            Cancel
          </Button>
        </Modal.Footer>
      </Modal>
    </div>
  );
}
