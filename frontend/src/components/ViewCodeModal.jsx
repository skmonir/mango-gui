import {PrismLight as SyntaxHighlighter} from 'react-syntax-highlighter';
import jsx from 'react-syntax-highlighter/dist/esm/languages/prism/jsx';
import {darcula} from 'react-syntax-highlighter/dist/esm/styles/prism';
import {useEffect, useState} from "react";
import {Button, Modal} from "react-bootstrap";
import DataService from "../services/DataService.js";

SyntaxHighlighter.registerLanguage('jsx', jsx);

export default function ViewCodeModal({codePath, setShowCodeModal}) {

    const [code, setCode] = useState('');
    const [showModal, setShowModal] = useState(false);

    useEffect(() => {
        console.log(codePath);
        console.log('ViewCodeModal is here');
        if (codePath) {
            fetchCode(codePath);
        }
    }, []);

    const fetchCode = (filepath) => {
        DataService.getCode({filePath: filepath}).then(code => {
            console.log(code);
            setCode(code);
            setShowModal(true);
        }).finally(() => setShowModal(true));
    }

    const closeModal = () => {
        setShowModal(false);
        setTimeout(() => setShowCodeModal(false), 500);
    }

    return (
        <Modal
            show={showModal}
            onHide={closeModal}
            size="lg"
            aria-labelledby="contained-modal-title-vcenter"
            centered
            fullscreen={true}
        >
            <Modal.Header/>
            <Modal.Body style={{height: '80vh', overflowY: 'auto'}}>
                <SyntaxHighlighter language="java" style={darcula}>
                    {code}
                </SyntaxHighlighter>
            </Modal.Body>
            <Modal.Footer>
                <Button size="sm" variant="outline-primary" onClick={() => closeModal()}>Close Code</Button>
            </Modal.Footer>
        </Modal>
    );
}