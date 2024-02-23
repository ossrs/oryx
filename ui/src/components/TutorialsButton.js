//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {Col, Row, Toast} from "react-bootstrap";
import logo from '../resources/logo.svg';
import * as Icon from 'react-bootstrap-icons';
import {Token} from "../utils";
import axios from "axios";
import Container from "react-bootstrap/Container";
import {useTranslation} from "react-i18next";
import {useSrsLanguage} from "./LanguageSwitch";

/**
 * Fetch the video tutorials from bilibili, for example:
 * @param bilibili The ref for video id of bilibili, must be a ref to avoid duplicated loading.
 * @returns A state of tutorials.
 */
function useTutorials({bilibili, medium}) {
  const language = useSrsLanguage();
  const cn = useTutorialsCn(bilibili);
  const en = useTutorialsEn(medium);
  return language === 'zh' ? cn : en;
}

function useTutorialsEn(mediumRef) {
  const language = useSrsLanguage();
  const [tutorials, setTutorials] = React.useState([]);
  const ref = React.useRef({tutorials:[]});

  const dict = React.useRef({
    '544e1db671c2': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/expand-your-global-reach-with-srs-stack-effortless-video-translation-and-dubbing-solutions-544e1db671c2',
      title: 'Revolutionize Video Content with SRS Stack: Effortless Dubbing and Translating to Multiple Languages Using OpenAI'
    },
    '13e28adf1e18': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/transform-your-browser-into-a-personal-voice-driven-gpt-ai-assistant-with-srs-stack-13e28adf1e18',
      title: 'Speak to the Future: Transform Your Browser into a Personal Voice-Driven GPT AI Assistant with SRS Stack'
    },
    'b3011e390e38': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/unlock-universal-ultra-low-latency-achieving-5-second-hls-live-streams-for-all-no-special-gear-b3011e390e38',
      title: 'Unlock Universal Ultra-Low Latency: Achieving 5-Second HLS Live Streams for All, No Special Equipment Needed'
    },
    'bb19c2a3bb7a': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/effortlessly-create-a-public-internet-whip-service-for-obs-a-guide-for-sub-second-streaming-bb19c2a3bb7a',
      title: 'Effortlessly Create a Public Internet WHIP Service for OBS: A Comprehensive Guide to Sub-Second Streaming'
    },
    '68PIGFDGihU': {
      author: 'Mr Bao',
      link: 'https://youtu.be/68PIGFDGihU',
      title: 'Live Stream: Simple, Budget-Friendly, No PC Needed! Ideal for Calm Media, Sleep Tunes, ASMR, Film Streaming, etc.'
    },
    'NBsdqUfKoOk': {
      author: 'Mr Bao',
      link: 'https://youtu.be/NBsdqUfKoOk',
      title: 'AI-Talk allows you to talk with OpenAI GPT.'
    },
    'nNOBFRshO6Q': {
      author: 'Mr Bao',
      link: 'https://youtu.be/nNOBFRshO6Q',
      title: '24/7 Live Stream: Easy Stream Your Camera to YouTube with DDNS & VPS - No PC or OBS Required!'
    },
    '1e902ab856bd': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/revolutionizing-live-streams-with-ai-transcription-creating-accessible-multilingual-subtitles-1e902ab856bd',
      title: 'Revolutionizing Live Streams with AI Transcription: Creating Accessible, Multilingual Subtitles for Diverse Audiences'
    },
    '39bd001af02d': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/efficient-live-streaming-transcoding-for-reducing-bandwidth-and-saving-costs-39bd001af02d',
      title: 'Efficient Live Streaming Transcoding for Reducing Bandwidth and Saving Costs'
    },
    'c078db917149': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/easily-stream-your-rtsp-ip-camera-to-youtube-twitch-or-facebook-c078db917149',
      title: 'Easily Stream Your RTSP IP Camera to YouTube, Twitch, or Facebook'
    },
    'ba1895828b4f': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/virtual-live-events-revolutionizing-the-way-we-experience-entertainment-ba1895828b4f',
      title: 'Mastering Virtual Live Events: Harness the Power of Pre-Recorded Content for Seamless and Engaging Live Streaming Experiences'
    },
    '2aa792c35b25': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/how-to-record-live-streaming-to-mp4-file-2aa792c35b25',
      title: 'Effortless Live Stream Recording with SRS Stack: A Step-by-Step Guide to Server-Side Recording and AWS S3 Integration'
    },
    '9748ae754c8c': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/how-to-setup-a-video-streaming-service-by-aapanel-9748ae754c8c',
      title: 'How to Setup a Video Streaming Service with aaPanel'
    },
    '38be22beec57': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/maximize-your-live-streaming-reach-a-guide-to-multi-platform-streaming-38be22beec57',
      title: 'Maximize Audience Engagement: Effortlessly Restream Live Content Across Multiple Platforms with SRS Stack'
    },
    'e9fe6f314ac6': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/how-to-setup-a-video-streaming-service-by-1-click-e9fe6f314ac6',
      title: 'How to Setup a Video Streaming Service by 1-Click',
    },
    'cb618777639f': {
      author: 'Winlin Yang',
      link: 'https://blog.ossrs.io/how-to-secure-srs-with-lets-encrypt-by-1-click-cb618777639f',
      title: 'How to Secure SRS Stack with Let’s Encrypt by 1-Click',
    },
    'ec18dfae7d6f': {
      author: 'Roboin',
      link: 'https://blog.ossrs.io/publish-your-srs-livestream-through-wordpress-ec18dfae7d6f',
      title: 'How to Publish Your SRS Livestream Through WordPress',
    },
  });

  const bvids = mediumRef?.current;
  React.useEffect(() => {
    if (!bvids || !bvids.length) return;
    if (language !== 'en') return;
    bvids.map(tutorial => {
      const obj = dict.current[tutorial.id];
      ref.current.tutorials.push({
        media: obj?.link?.indexOf('youtu.be') > 0 ? 'YouTube' : 'Medium',
        ...obj,
      });
      setTutorials([...ref.current.tutorials]);
      return null;
    });
  }, [bvids, language]);

  return tutorials;
}

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

    bvids.map(tutorial => {
      tutorial.link = `https://www.bilibili.com/video/${tutorial.id}`;

      axios.post(`/terraform/v1/mgmt/bilibili`, {
        bvid: tutorial.id,
      }, {
        headers: Token.loadBearerHeader(),
        cancelToken: source.token,
      }).then(res => {
        const data = res.data.data;
        tutorial.media = 'Bilibili';
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
    <Container fluid>
      <Row>
        {tutorials.map((tutorial, index) => (
          <Col key={index} sm={3}>
            <Toast onClose={onClose}>
              <Toast.Header>
                <img src={logo} className="rounded me-2" width={56} alt=''/>
                <strong className="me-auto">{tutorial.media}</strong>
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

