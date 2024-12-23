package controller

import (
	"context"
	"os/exec"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/adeptvin1/kubernetes-operator-for-LWMECPS/api/v1alpha1" // Замените на правильный путь
)

// LitmusInstallReconciler контроллер для установки Litmus Chaos Operator
type LitmusInstallReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile проверяет и устанавливает Litmus Chaos Operator, если его нет
func (r *LitmusInstallReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("litmus_install", req.NamespacedName)

	// Проверяем, установлен ли Litmus Chaos Operator (проверка CRD)
	if err := ensureLitmusInstalled(ctx, log); err != nil {
		log.Error(err, "Не удалось установить Litmus Chaos Operator")
		return ctrl.Result{}, err
	}

	log.Info("Litmus Chaos Operator успешно установлен")
	return ctrl.Result{}, nil
}

// ensureLitmusInstalled проверяет наличие Litmus и устанавливает его через Helm, если необходимо
func ensureLitmusInstalled(ctx context.Context, log logr.Logger) error {
	log.Info("Проверка наличия Litmus Chaos Operator")

	// Проверяем наличие CRD Litmus (например, ChaosEngine)
	requiredCRD := "chaosengines.litmuschaos.io"
	if !checkCRD(ctx, requiredCRD) {
		log.Info("Litmus Chaos Operator не найден, выполняется установка...")

		// Устанавливаем Litmus через Helm
		cmd := exec.Command("helm", "install", "litmus", "litmuschaos/litmus",
			"--namespace", "litmus", "--create-namespace")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(err, "Не удалось установить Litmus Chaos Operator", "output", string(output))
			return err
		}

		log.Info("Litmus Chaos Operator успешно установлен")
	} else {
		log.Info("Litmus Chaos Operator уже установлен")
	}

	return nil
}

// checkCRD проверяет, существует ли CRD с заданным именем
func checkCRD(ctx context.Context, crdName string) bool {
	cmd := exec.Command("kubectl", "get", "crd", crdName)
	err := cmd.Run()
	return err == nil
}

// SetupWithManager регистрирует этот контроллер с контроллером-менеджером
func (r *LitmusInstallReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Замените на ваш правильный тип CRD, например LitmusInstall
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.LitmusInstall{}). // Используем тип LitmusInstall из API
		Complete(r)
}
