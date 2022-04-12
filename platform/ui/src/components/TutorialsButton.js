import React from "react";
import {Col, Row, Toast} from "react-bootstrap";
import logo from '../resources/logo.svg';
import * as Icon from 'react-bootstrap-icons';
import {Token} from "../utils";
import axios from "axios";
import Container from "react-bootstrap/Container";
import {useTranslation} from "react-i18next";
import {useSrsLanguage} from "./LanguageSwitch";

function useTutorials(bvidsRef, mediumRef) {
  const language = useSrsLanguage();
  const cn = useTutorialsCn(bvidsRef);
  const en = useTutorialsEn(mediumRef);
  return language === 'zh' ? cn : en;
}

function useTutorialsEn(mediumRef) {
  const language = useSrsLanguage();
  const [tutorials, setTutorials] = React.useState([]);
  const ref = React.useRef({tutorials:[]});

  const dict = React.useRef({
    'e9fe6f314ac6': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/how-to-setup-a-video-streaming-service-by-1-click-e9fe6f314ac6',
      title: 'How to Setup a Video Streaming Service by 1-Click',
    },
    'cb618777639f': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/how-to-secure-srs-with-lets-encrypt-by-1-click-cb618777639f',
      title: 'How to Secure SRS Droplet with Let’s Encrypt by 1-Click',
    },
  });

  const bvids = mediumRef?.current;
  React.useEffect(() => {
    if (!bvids || !bvids.length) return;
    if (language !== 'en') return;
    bvids.map(tutorial => {
      ref.current.tutorials.push(dict.current[tutorial.id]);
      setTutorials([...ref.current.tutorials].sort((a, b) => b.view - a.view));
      return null;
    });
  }, [bvids, language]);

  return tutorials;
}

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
function useTutorialsCn(bvidsRef) {
  const language = useSrsLanguage();
  const [tutorials, setTutorials] = React.useState([]);
  const ref = React.useRef({tutorials:[]});

  const bvids = bvidsRef.current;
  React.useEffect(() => {
    if (!bvids || !bvids.length) return;
    if (language !== 'zh') return;

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
  }, [bvids, language]);

  return tutorials;
}

// A toast list for tutorials.
function TutorialsToast({tutorials, onClose}) {
  const {t} = useTranslation();

  return (<>
    <Container>
      <Row>
        {tutorials.map((tutorial, index) => (
          <Col xs lg={4} key={index}>
            <Toast onClose={onClose}>
              <Toast.Header>
                <img src={logo} className="rounded me-2" width={56} alt=''/>
                <strong className="me-auto">{t('tutorials.media')}</strong>
                {tutorial.view && <> <span title={t('tutorials.view')}><Icon.Play /> {tutorial.view}</span> &nbsp; </>}
                {tutorial.like && <> <span title={t('tutorials.like')}><Icon.HandThumbsUp /> {tutorial.like}</span> &nbsp; </>}
                {tutorial.share && <> <span title={t('tutorials.share')}><Icon.Share /> {tutorial.share}</span> &nbsp; </>}
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

