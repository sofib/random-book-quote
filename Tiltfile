load('ext://restart_process', 'docker_build_with_restart')
load('ext://helm_resource', 'helm_resource', 'helm_repo')

tools_path = os.environ.get('TOOLS_PATH', os.path.join(os.getcwd(), 'tools'))
def tools(cmd):
    return os.path.join(tools_path, cmd)

IMAGE_NAME = 'random-quote'
NAMESPACE = 'random-quote'
REGISTRY_HOST = 'localhost:5005'

allow_k8s_contexts('kind-random-quote-cluster')

local_resource(
    'create-cluster',
    cmd=tools('ctlptl apply -f kind/cluster.yaml'),
    auto_init=True,
    trigger_mode=TRIGGER_MODE_AUTO,
    labels=['cluster']
)

local_resource(
    'create-registry',
    cmd=tools('ctlptl apply -f kind/registry.yaml'),
    auto_init=True,
    trigger_mode=TRIGGER_MODE_AUTO,
    resource_deps=['create-cluster'],
    labels=['cluster']
)

local_resource(
    'destroy-cluster',
    cmd=tools('ctlptl delete -f kind/cluster.yaml'),
    auto_init=False,
    trigger_mode=TRIGGER_MODE_MANUAL,
    resource_deps=['create-cluster'],
    labels=['cluster']
)

local_resource(
    'create-namespace',
    cmd=tools('kubectl apply -f kind/namespace.yaml'),
    resource_deps=['create-cluster'],
    trigger_mode=TRIGGER_MODE_AUTO,
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
    labels=['random-quote'],
    resource_deps=['create-namespace'],
)

local_resource(
    'manual-run',
    cmd=tools('kubectl create job --from cronjob/random-quote -n %s test-job-$(date +%%s)' % NAMESPACE),
    labels=['random-quote'],
    resource_deps=['create-namespace'],
    auto_init=False,
    trigger_mode=TRIGGER_MODE_MANUAL,
)