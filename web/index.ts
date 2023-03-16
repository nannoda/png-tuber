import { REMOTE_ID_API } from "./consts";

async function getText(url: string) {
    const response = await fetch(url);
    return await response.text();
}
console.log("3");

async function postRemoteId(localId: string) {

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
        console.log(`type: ${typeof event.data}`);

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
                const localId = btoa(JSON.stringify(peerConnection.localDescription));
                console.log("localDescription: " + localId);
                postRemoteId(localId);
            }
        }
    }

    const offer = await peerConnection.createOffer();
    await peerConnection.setLocalDescription(offer);

}

main();