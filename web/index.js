(() => {
  // consts.ts
  var REMOTE_ID_API = "/api/remoteId";

  // index.ts
  console.log("3");
  async function postRemoteId(localId) {
    const response = await fetch(REMOTE_ID_API, {
      method: "POST",
      body: localId
    });
    console.log(response);
  }
  async function main() {
    const peerConnection = new RTCPeerConnection({
      iceServers: [
        {
          urls: "stun:stun.l.google.com:19302"
        }
      ]
    });
    const setSession = (base64Str) => {
      const session = JSON.parse(atob(base64Str));
      peerConnection.setRemoteDescription(session);
    };
    window.setSession = setSession;
    const sendChannel = peerConnection.createDataChannel("sendChannel");
    sendChannel.onopen = () => {
      console.log("sendChannel open");
    };
    sendChannel.onmessage = (event) => {
      console.log(`type: ${typeof event.data}`);
      console.log("sendChannel message: " + event.data);
    };
    sendChannel.onclose = () => {
      console.log("sendChannel close");
    };
    peerConnection.onicecandidate = (event) => {
      if (event.candidate) {
        console.log("candidate: " + event.candidate.candidate);
      } else {
        console.log("candidate: null");
        if (peerConnection.localDescription) {
          const localId = btoa(JSON.stringify(peerConnection.localDescription));
          console.log("localDescription: " + localId);
          postRemoteId(localId);
        }
      }
    };
    const offer = await peerConnection.createOffer();
    await peerConnection.setLocalDescription(offer);
  }
  main();
})();
//# sourceMappingURL=index.js.map
