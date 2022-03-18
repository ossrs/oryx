import React from "react";
import {Col, Row, Toast} from "react-bootstrap";
import logo from '../resources/logo.svg';
import * as Icon from 'react-bootstrap-icons';
import {Token} from "../utils";
import axios from "axios";
import Container from "react-bootstrap/Container";

/**
 * Fetch the video tutorials from bilibili, for example:
 * @param bvidsRef The ref for video id of bilibili, must be a ref to avoid duplicated loading.
 * @returns A state of tutorials.
 *
 * For example:
      const sslTutorials = useTutorials(React.useRef([
        {author: '程晓龙', id: 'BV1tZ4y1R7qp'},
      ]));
 * Then, use the state in TutorialsButton:
      return <TutorialsButton prefixLine={true} tutorials={sslTutorials} />
 */
function useTutorials(bvidsRef) {
  const bvids = bvidsRef.current;

  const [tutorials, setTutorials] = React.useState([]);
  const ref = React.useRef({tutorials:[]});

  React.useEffect(() => {
    if (!bvids || !bvids.length) return;

    // Allow cancel up the requests.
    const source = axios.CancelToken.source();

    const token = Token.load();
    bvids.map(tutorial => {
      tutorial.link = `https://www.bilibili.com/video/${tutorial.id}`;

      axios.post(`/terraform/v1/mgmt/bilibili`, {
        ...token, bvid: tutorial.id,
      }, {
        cancelToken: source.token,
      }).then(res => {
        const data = res.data.data;
        tutorial.title = data.title;
        tutorial.desc = data.desc;
        tutorial.view = parseInt(data.stat.view);
        tutorial.like = parseInt(data.stat.like);
        tutorial.share = parseInt(data.stat.share);
        // Order by view desc.
        ref.current.tutorials.push(tutorial);
        setTutorials([...ref.current.tutorials].sort((a, b) => b.view - a.view));
      }).catch((e) => {
        if (axios.isCancel(e)) return;
        throw e;
      });
      return null;
    });

    return () => {
      // When cleanup, cancel all requests to avoid update the unmounted components, like error message as:
      //    Can't perform a React state update on an unmounted component.
      //    This is a no-op, but it indicates a memory leak in your application.
      source.cancel();
    };
  }, [bvids]);

  return tutorials;
}

// A toast list for tutorials.
function TutorialsToast({tutorials, onClose}) {
  return (<>
    <Container>
      <Row>
        {tutorials.map((tutorial, index) => (
          <Col xs lg={4} key={index}>
            <Toast onClose={onClose}>
              <Toast.Header>
                <img src={logo} className="rounded me-2" width={56} alt=''/>
                <strong className="me-auto">Bilibili</strong>
                <span title='播放次数'><Icon.Play /> {tutorial.view}</span> &nbsp;
                <span title='点赞次数'><Icon.HandThumbsUp /> {tutorial.like}</span> &nbsp;
                <span title='分享次数'><Icon.Share /> {tutorial.share}</span> &nbsp;
                <small>by {tutorial.author}</small>
              </Toast.Header>
              <Toast.Body>
                <a href={tutorial.link} target='_blank' rel='noreferrer'>
                  {tutorial.title}
                </a>
              </Toast.Body>
            </Toast>
            <p></p>
          </Col>
        ))}
      </Row>
    </Container>
  </>);
}

// The tutorials button, the props tutorials is a array, create by useTutorials.
function TutorialsButton({tutorials, prefixLine}) {
  const [show, setShow] = React.useState(false);

  return (
    <>
      <div role='button' style={{display: 'inline-block'}}>
        <Icon.PatchQuestion onClick={() => setShow(!show)} />
      </div>
      {show && prefixLine && <p></p>}
      {show &&
        <TutorialsToast
          prefixLine={prefixLine}
          tutorials={tutorials}
          onClose={() => setShow(false)}
        />
      }
    </>
  );
}

export {useTutorials, TutorialsButton, TutorialsToast};

