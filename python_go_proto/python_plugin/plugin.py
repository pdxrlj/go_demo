import sys
import time
import grpc
from concurrent import futures
import hello_pb2_grpc
import hello_pb2

from grpc_health.v1.health import HealthServicer
from grpc_health.v1 import health_pb2, health_pb2_grpc


class GreetServer(hello_pb2_grpc.GreetingServiceServicer):
    def Greet(self, request, context):
        response = hello_pb2.GreetingResponse()  # _GREETINGRESPONSE
        response.message = f"Hello {request.name}!"
        return response


def serve():
    health = HealthServicer()

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    health.set("plugin", health_pb2.HealthCheckResponse.ServingStatus.Value('SERVING'))
    health_pb2_grpc.add_HealthServicer_to_server(GreetServer(), server)

    server.add_insecure_port(':1234')
    server.start()
    # go plugin will read this line
    print("1|1|tcp|127.0.0.1:1234|grpc")
    sys.stdout.flush()

    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == "__main__":
    serve()
