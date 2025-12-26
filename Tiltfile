load('ext://restart_process', 'docker_build_with_restart')
load('ext://helm_resource', 'helm_resource', 'helm_repo')

tools_path = os.environ.get('TOOLS_PATH', os.path.join(os.getcwd(), 'tools'))
def tools(cmd):
    return os.path.join(tools_path, cmd)

IMAGE_NAME = 'random-quote'
NAMESPACE = 'random-quote'
REGISTRY_HOST = 'localhost:5005'

allow_k8s_contexts('kind-random-quote-cluster')

print('Ensuring base cluster setup exist...')
local(tools('ctlptl apply -f kind/cluster.yaml'), quiet=False, echo_off=False)
local(tools('ctlptl apply -f kind/registry.yaml'), quiet=False, echo_off=False)
local(tools('kubectl apply -f kind/namespace.yaml'), quiet=False, echo_off=False)

print(local(tools('kubectl config current-context'), quiet=False, echo_off=False))
print(local('cat ~/.kube/config', quiet=False, echo_off=False))

print('Cluster setup complete')

k8s_context('kind-random-quote-cluster')

local_resource(
    'destroy-cluster',
    cmd=tools('ctlptl delete -f kind/cluster.yaml'),
    auto_init=False,
    trigger_mode=TRIGGER_MODE_MANUAL,
    labels=['cluster']
)

docker_build(
    '%s/%s' % (REGISTRY_HOST, IMAGE_NAME),
    '.',
    dockerfile='Dockerfile',
    build_args={
        'BASE_IMAGE': 'alpine:latest',
    },
    only=[
        'go.mod',
        'go.sum',
        'main.go',
        'internal/',
    ],
    live_update=[
        sync('./internal', '/internal'),
        run('go build -o random-quote main.go', trigger=['./internal'])
    ]
)

helm_resource(
    'random-quote',
    './helm',
    namespace=NAMESPACE,
    flags=[
        '--set', 'image.repository=%s/%s' % (REGISTRY_HOST, IMAGE_NAME),
        '--set', 'image.pullPolicy=Always',
    ],
    image_deps=['%s/%s' % (REGISTRY_HOST, IMAGE_NAME)],
    image_keys=[('image.repository', 'image.tag')],
    labels=['random-quote']
)

local_resource(
    'manual-run',
    cmd=tools('kubectl create job --from cronjob/random-quote -n %s test-job-$(date +%%s)' % NAMESPACE),
    labels=['random-quote'],
    auto_init=False,
    trigger_mode=TRIGGER_MODE_MANUAL,
)