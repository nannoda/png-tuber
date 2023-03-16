(() => {
  // index.ts
  console.log("2");
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
          console.log("localDescription: " + btoa(JSON.stringify(peerConnection.localDescription)));
        }
      }
    };
    const offer = await peerConnection.createOffer();
    await peerConnection.setLocalDescription(offer);
  }
  main();
})();
//# sourceMappingURL=index.js.map
