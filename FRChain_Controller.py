from ryu.base import app_manager
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.ofproto import ofproto_v1_3
from ryu.lib.packet import packet, ethernet
import requests

BLOCKCHAIN_API = "http://blockchain-network/api/v1/validate"

class FRChainController(app_manager.RyuApp):
    OFP_VERSIONS = [ofproto_v1_3.OFP_VERSION]

    def __init__(self, *args, **kwargs):
        super(FRChainController, self).__init__(*args, **kwargs)

    def validate_flow(self, flow_rule):
        # Send flow rule for validation to the blockchain
        response = requests.post(BLOCKCHAIN_API, json=flow_rule)
        if response.status_code == 200:
            return response.json().get("valid", False)
        return False

    @set_ev_cls(ofp_event.EventOFPPacketIn, MAIN_DISPATCHER)
    def packet_in_handler(self, ev):
        msg = ev.msg
        datapath = msg.datapath
        ofproto = datapath.ofproto
        parser = datapath.ofproto_parser
        pkt = packet.Packet(msg.data)
        eth = pkt.get_protocol(ethernet.ethernet)

        # Create a flow rule based on the packet
        in_port = msg.match['in_port']
        flow_rule = {
            "src": eth.src,
            "dst": eth.dst,
            "in_port": in_port,
            "controller": "controller_1"
        }

        # Validate flow rule
        if not self.validate_flow(flow_rule):
            self.logger.warning("Flow rule rejected by blockchain!")
            return

        # If valid, install the flow
        self.logger.info("Flow rule validated and installed.")
        actions = [parser.OFPActionOutput(ofproto.OFPP_FLOOD)]
        out = parser.OFPPacketOut(datapath=datapath, buffer_id=msg.buffer_id,
                                  in_port=in_port, actions=actions, data=msg.data)
        datapath.send_msg(out)
