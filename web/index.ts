async function getText(url: string) {
    const response = await fetch(url);
    return await response.text();
}
console.log("2");



async function main() {
    const peerConnection = new RTCPeerConnection({
        iceServers: [
            {
                urls: "stun:stun.l.google.com:19302"
            }
        ]
    });

    const setSession = (base64Str: string) => {
        const session = JSON.parse(atob(base64Str));
        peerConnection.setRemoteDescription(session);
    }
    (window as any).setSession = setSession;

    const sendChannel = peerConnection.createDataChannel("sendChannel");
    sendChannel.onopen = () => {
        console.log("sendChannel open");
    }

    sendChannel.onmessage = (event) => {
        console.log("sendChannel message: " + event.data);
    }

    sendChannel.onclose = () => {
        console.log("sendChannel close");
    }

    peerConnection.onicecandidate = (event) => {
        if (event.candidate) {
            console.log("candidate: " + event.candidate.candidate);
        } else {
            console.log("candidate: null");
            if (peerConnection.localDescription) {
                console.log("localDescription: " + btoa(JSON.stringify(peerConnection.localDescription)));
            }
        }
    }

    const offer = await peerConnection.createOffer();
    await peerConnection.setLocalDescription(offer);
}

main();